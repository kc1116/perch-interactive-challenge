package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

var logger = logrus.New()

const (
	mqttClientIDFMT      = "projects/%s/locations/%s/registries/%s/devices/%s"
	interactionsTopicFMT = "/devices/%s/events/interactions"
	parentFMT            = "projects/%s/locations/%s"
)

func Logger() *logrus.Logger {
	return logger
}

// MQTTClientID
func MQTTClientID(projectID, region, registryID, deviceID string) string {
	return fmt.Sprintf(mqttClientIDFMT, projectID, region, registryID, deviceID)
}

// DeviceRegistryParent
func DeviceRegistryParent(projectID, region string) string {
	return fmt.Sprintf(parentFMT, projectID, region)
}

// https://ericchiang.github.io/post/go-tls/ this blog post was the only post that properly explained TLS Certs etc
// before that google iot-core docs are super vague about what you need to actually do

// helper function to create a cert template with a serial number and other required fields
func CertTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{CommonName: "unused"},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func CreateCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}

	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func CreateKey() (*rsa.PrivateKey, *x509.Certificate, error) {
	// generate a new key-pair
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("generating random key: %v", err)
	}

	rootCertTmpl, err := CertTemplate()
	if err != nil {
		return nil, nil, fmt.Errorf("creating cert template: %v", err)
	}

	// describe what the certificate will be used for
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}

	return rootKey, rootCertTmpl, nil
}

// EncodeEvent base64 encode event proto bytes
func EncodeEvent(evt *protos.Event) (string, error) {
	b, err := proto.Marshal(evt)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// DecodeEvt decode incoming event
func DecodeEvt(encodedEvtStr string) *protos.Event {
	evt := &protos.Event{}
	b, _ := hex.DecodeString(encodedEvtStr)

	_ = proto.Unmarshal(b, evt)
	return evt
}
