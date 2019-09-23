package core

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/cloudiot/v1"
	"strings"
)

type DeviceRegistry struct {
	Region       string
	RegistryID   string
	TopicID      string
	projectID    string
	parent       string
	Registry     *cloudiot.DeviceRegistry
	Topic        *pubsub.Topic
	Client       *cloudiot.Service
	PubSubClient *pubsub.Client
}

// Init initializes our DeviceRegistry by creating the device Registry and pub/sub Topic in google cloud
func (d *DeviceRegistry) Init(create bool) (*DeviceRegistry, error) {
	gcclient, err := GCHttpClient()
	if err != nil {
		return nil, err
	}

	pubsubClient, err := PubSubClient(d.projectID)
	if err != nil {
		return nil, err
	}

	d.Client = gcclient
	d.PubSubClient = pubsubClient

	err = d.CreateTopic()
	if err != nil {
		return d, err
	}

	if registry := d.GetRegistry(); registry != nil {
		d.Registry = registry
		return d, nil
	}

	if !create {
		return nil, fmt.Errorf("Registry not found %s, if you want to create the Registry pass true for create", d.RegistryName())
	}

	registry, err := d.CreateRegistry(d.Topic.String())
	if err != nil {
		return d, err
	}

	d.Registry = registry
	return d, nil
}

// CreateRegistry creates a device Registry if it does not already exists
func (d *DeviceRegistry) CreateRegistry(fullTopicPath string) (*cloudiot.DeviceRegistry, error) {
	registry := &cloudiot.DeviceRegistry{
		Id: d.RegistryID,
		EventNotificationConfigs: []*cloudiot.EventNotificationConfig{
			{
				PubsubTopicName: fullTopicPath,
			},
		},
	}

	_, err := d.Client.Projects.Locations.Registries.Create(d.parent, registry).Do()
	if err != nil {
		return registry, err
	}

	d.Registry = registry

	return registry, err
}

// GetRegistry gets Registry details from GCP
func (d *DeviceRegistry) GetRegistry() *cloudiot.DeviceRegistry {
	if registry, err := d.Client.Projects.Locations.Registries.Get(d.RegistryName()).Do(); err == nil {
		return registry
	}

	return nil
}

// CreateTopic creates a Topic if it does not already exists
func (d *DeviceRegistry) CreateTopic() error {
	topic := d.PubSubClient.Topic(d.TopicID)
	if ok, _ := topic.Exists(context.Background()); ok {
		d.Topic = topic
		return nil
	}

	topic, err := d.PubSubClient.CreateTopic(context.Background(), d.TopicID)
	if err != nil {
		return err
	}

	d.Topic = topic
	return nil
}

// CleanUp destroys Registry resource in GCP
func (d *DeviceRegistry) CleanUp() error {
	if registry := d.GetRegistry(); registry != nil {
		_, err := d.Client.Projects.Locations.Registries.Delete(d.RegistryName()).Do()
		if err != nil {
			return err
		}
	}
	return nil
}

// RegistryName
func (d *DeviceRegistry) RegistryName() string {
	return fmt.Sprintf("projects/%s/locations/%s/registries/%s", d.projectID, d.Region, d.RegistryID)
}

// String
func (d *DeviceRegistry) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Region: %s \n", d.Region))
	builder.WriteString(fmt.Sprintf("RegistryID: %s \n", d.RegistryID))
	builder.WriteString(fmt.Sprintf("RegistryName: %s \n", d.Registry.Name))
	builder.WriteString(fmt.Sprintf("TopicID: %s \n", d.TopicID))
	builder.WriteString(fmt.Sprintf("ProjectID: %s \n", d.projectID))

	return builder.String()
}

// NewDeviceRegistry returns uninitialized DeviceRegistry struct
func NewDeviceRegistry(projectID, region, registryID, topicID string) *DeviceRegistry {
	return &DeviceRegistry{
		projectID:  projectID,
		RegistryID: registryID,
		Region:     region,
		TopicID:    topicID,
		parent:     DeviceRegistryParent(projectID, region),
	}
}
