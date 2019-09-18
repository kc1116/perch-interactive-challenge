package core

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/api/cloudiot/v1"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var certPool *x509.CertPool

const (
	mqttServer = "ssl://mqtt.googleapis.com:8883"
	topicType = "events"
	qos = 1
	retain = false
	username = "unused"
)

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
	Certs TLSCerts
}

type TLSCerts struct {
	RootCertTmpl *x509.Certificate
	Cert *x509.Certificate
	Pem string
	Key *rsa.PrivateKey
}

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
	if d.device = d.GetDevice(); d.device != nil {
		return nil
	}

	err = d.NewKey()
	if err != nil {
		return err
	}

	err = d.JWT()
	if err != nil {
		return err
	}

	log.Println(d.Certs.Pem)
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

func (d *Device) ConnectMQTT() error {
	config := &tls.Config{
		RootCAs:            certPool,
		MinVersion: tls.VersionTLS12,
	}

	opts := mqtt.NewClientOptions().
		AddBroker(mqttServer).
		SetTLSConfig(config).
		SetClientID(d.DevicePath).
		SetUsername(username).
		SetPassword(d.tokenString).
		SetProtocolVersion(4)

	log.Println(opts.ClientID)
	d.mqttconn = mqtt.NewClient(opts)
	if token := d.mqttconn.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Failed to connect client ", token.Error())
		return token.Error()
	}

	return nil
}

func (d *Device) NewKey() error {
	key, rootCertTempl, err := CreateKey()
	rootCert, rootCertPEM, err := CreateCert(rootCertTempl, rootCertTempl, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("error creating cert: %v", err)
	}

	d.Certs = TLSCerts{
		RootCertTmpl: rootCertTempl,
		Cert: rootCert,
		Pem: fmt.Sprintf("%s", rootCertPEM),
		Key: key,
	}

	return nil
}

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

// String
func (d *Device) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DeviceID: %s \n", d.device.Id))
	builder.WriteString(fmt.Sprintf("DeviceName: %s \n", d.device.Name))
	builder.WriteString(fmt.Sprintf("Region: %s \n", d.Region))
	builder.WriteString(fmt.Sprintf("RegistryID: %s \n", d.RegistryID))
	builder.WriteString(fmt.Sprintf("ProjectID: %s \n", d.projectID))

	return builder.String()
}

func DeviceID() string {
	return fmt.Sprintf("perchdevice-%04X%04X", rand.Intn(0x10000), rand.Intn(0x10000))
}

// NewDevice returns unintialized device struct
func NewDevice( projectID, region, registryID, registryPath string) *Device {
	deviceID :=  DeviceID()
	return &Device{
		Region: region,
		projectID: projectID,
		RegistryID: registryID,
		DeviceID:  deviceID,
		parent: registryPath,
		DevicePath: fmt.Sprintf("%s/devices/%s", registryPath, deviceID),
		eventTopic: fmt.Sprintf("/devices/%s/events", deviceID),
	}
}

func init()  {
	log.Println("creating cert pool for MQTT connections")
	certPool = x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("/Users/kc1116/Desktop/perch-interactive-challenge/core/google-cert/roots.pem")
	if err != nil {
		log.Fatal("can not load google certs ", err)
	}

	certPool.AppendCertsFromPEM(pemCerts)

	mqtt.DEBUG = log.New(os.Stderr, "DEBUG - ", log.LstdFlags)
	mqtt.CRITICAL = log.New(os.Stderr, "CRITICAL - ", log.LstdFlags)
	mqtt.WARN = log.New(os.Stderr, "WARN - ", log.LstdFlags)
	mqtt.ERROR = log.New(os.Stderr, "ERROR - ", log.LstdFlags)
}