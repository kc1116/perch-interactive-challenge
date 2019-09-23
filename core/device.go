package core

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	"google.golang.org/api/cloudiot/v1"
	"io/ioutil"
	"math/rand"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	mqttServer = "ssl://mqtt.googleapis.com:8883"
	topicType  = "events"
	qos        = 1
	retain     = false
	username   = "unused"
)

var mqttPool = &MqttPool{}

type MqttPool struct {
	certs *x509.CertPool
	conn  mqtt.Client
	sync.Once
	sync.Mutex
}

type Device struct {
	Region      string
	DeviceID    string
	RegistryID  string
	DevicePath  string
	projectID   string
	parent      string
	eventTopic  string
	tokenString string
	device      *cloudiot.Device
	client      *cloudiot.Service
	mqttconn    mqtt.Client
	Certs       TLSCerts
}

type TLSCerts struct {
	RootCertTmpl *x509.Certificate
	Cert         *x509.Certificate
	Pem          string
	Key          *rsa.PrivateKey
}

// Init
func (d *Device) Init() (*Device, error) {
	var err error

	d.client, err = GCHttpClient()
	if err != nil {
		return nil, err
	}

	err = d.CreateDevice()

	return d, err
}

// CreateDevice creates our device in google cloud
func (d *Device) CreateDevice() error {
	var err error
	_ = d.CleanUp()

	err = d.NewKey()
	if err != nil {
		return err
	}

	err = d.JWT()
	if err != nil {
		return err
	}

	device := cloudiot.Device{
		Id: d.DeviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "RSA_X509_PEM",
					Key:    d.Certs.Pem,
				},
			},
		},
	}

	d.device, err = d.client.Projects.Locations.Registries.Devices.Create(d.parent, &device).Do()
	if err != nil {
		return err
	}

	return nil
}

// GetDevice
func (d *Device) GetDevice() *cloudiot.Device {
	if device, err := d.client.Projects.Locations.Registries.Devices.Get(d.DevicePath).Do(); err == nil {
		return device
	}

	return nil
}

// CleanUp
func (d *Device) CleanUp() error {
	if device := d.GetDevice(); device != nil {
		_, err := d.client.Projects.Locations.Registries.Devices.Delete(d.DevicePath).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

// NewKey
func (d *Device) NewKey() error {
	key, rootCertTempl, err := CreateKey()
	rootCert, rootCertPEM, err := CreateCert(rootCertTempl, rootCertTempl, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("error creating cert: %v", err)
	}

	d.Certs = TLSCerts{
		RootCertTmpl: rootCertTempl,
		Cert:         rootCert,
		Pem:          fmt.Sprintf("%s", rootCertPEM),
		Key:          key,
	}

	return nil
}

// JWT
func (d *Device) JWT() error {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		Audience:  d.projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	tokenString, err := token.SignedString(d.Certs.Key)
	if err != nil {
		return err
	}

	d.tokenString = tokenString
	return nil
}

// ConnectMQTT
func (d *Device) ConnectMQTT() error {
	var err error
	mqttPool.Do(func() {
		mqttPool.Lock()

		config := &tls.Config{
			RootCAs:    mqttPool.certs,
			MinVersion: tls.VersionTLS12,
		}

		opts := mqtt.NewClientOptions().
			AddBroker(mqttServer).
			SetTLSConfig(config).
			SetClientID(d.DevicePath).
			SetUsername(username).
			SetPassword(d.tokenString).
			SetProtocolVersion(4)

		mqttPool.conn = mqtt.NewClient(opts)
		if token := mqttPool.conn.Connect(); token.Wait() && token.Error() != nil {
			err = token.Error()
		}

		mqttPool.Unlock()
	})

	return err
}

// Publish
func (d *Device) Publish(evt *protos.Event) error {
	encoded, err := EncodeEvent(evt)
	if err != nil {
		return fmt.Errorf("error encoding evt during publish %s", err)
	}
	topic := fmt.Sprintf(interactionsTopicFMT, d.device.Id)

	mqttPool.Lock()
	if token := mqttPool.conn.Publish(topic, qos, retain, encoded); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	mqttPool.Unlock()

	return nil
}

// StartSession
func (d *Device) StartSession(wg *sync.WaitGroup) {
	NewSession(d.DeviceID, d.Publish).Start(wg)
}

// String
func (d *Device) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ID: %s \n", d.device.Id))
	builder.WriteString(fmt.Sprintf("DeviceName: %s \n", d.device.Name))
	builder.WriteString(fmt.Sprintf("Region: %s \n", d.Region))
	builder.WriteString(fmt.Sprintf("RegistryID: %s \n", d.RegistryID))
	builder.WriteString(fmt.Sprintf("ProjectID: %s \n", d.projectID))

	return builder.String()
}

// ID
func ID() string {
	return fmt.Sprintf("perchdevice-%04X%04X", rand.Intn(0x10000), rand.Intn(0x10000))
}

// NewDevice returns unintialized device struct
func NewDevice(projectID, region, registryID, registryPath string) *Device {
	deviceID := ID()
	return &Device{
		Region:     region,
		projectID:  projectID,
		RegistryID: registryID,
		DeviceID:   deviceID,
		parent:     registryPath,
		DevicePath: fmt.Sprintf("%s/devices/%s", registryPath, deviceID),
		eventTopic: fmt.Sprintf("/devices/%s/events", deviceID),
	}
}

func init() {
	mqttPool.certs = x509.NewCertPool()
	_, name, _, _ := runtime.Caller(0)
	p := path.Join(path.Dir(name), "./google-cert/roots.pem")

	pemCerts, err := ioutil.ReadFile(p)
	if err != nil {
		logger.Fatalf("can not load google certs %s", err)
	}

	mqttPool.certs.AppendCertsFromPEM(pemCerts)
}
