package core

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/cloudiot/v1"
	"strings"
	"sync"
)

type DeviceRegistry struct {
	Region       string
	RegistryID   string
	TopicID      string
	projectID    string
	parent string
	once sync.Once
	registry     *cloudiot.DeviceRegistry
	topic *pubsub.Topic
	client       *cloudiot.Service
	pubSubClient *pubsub.Client
}

// Init initializes our DeviceRegistry by creating the device registry and pub/sub topic in google cloud
func (d *DeviceRegistry) Init() (*DeviceRegistry, error) {
	gcclient, err := GCHttpClient()
	if err != nil {
		return nil, err
	}

	pubsubClient, err := PubSubClient(d.projectID)
	if err != nil {
		return nil, err
	}

	d.client = gcclient
	d.pubSubClient = pubsubClient

	if registry := d.GetRegistry(); registry != nil {
		d.registry = registry
		return d, nil
	}

	err = d.CreateTopic()
	if err != nil {
		return d, err
	}

	registry, err := d.CreateRegistry(d.topic.String())
	if err != nil {
		return d, err
	}

	d.registry = registry
	return d, nil
}

// CreateRegistry creates a device registry if it does not already exists
func (d *DeviceRegistry) CreateRegistry(fullTopicPath string) (*cloudiot.DeviceRegistry, error) {
	registry := &cloudiot.DeviceRegistry{
		Id: d.RegistryID,
		EventNotificationConfigs: []*cloudiot.EventNotificationConfig{
			{
				PubsubTopicName: fullTopicPath,
			},
		},
	}

	_, err := d.client.Projects.Locations.Registries.Create(d.parent, registry).Do()
	if err != nil {
		return registry, err
	}

	d.registry = registry

	return registry, err
}

func (d *DeviceRegistry) GetRegistry() *cloudiot.DeviceRegistry {
	if registry, err := d.client.Projects.Locations.Registries.Get(d.RegistryName()).Do(); err == nil {
		return registry
	}

	return nil
}

// CreateTopic creates a topic if it does not already exists
func (d *DeviceRegistry) CreateTopic() error {
	topic := d.pubSubClient.TopicInProject(d.TopicID, d.projectID)
	if topic != nil {
		d.topic = topic
		return nil
	}

	var err error
	topic, err = d.pubSubClient.CreateTopic(context.Background(), d.TopicID)
	if err != nil {
		return nil
	}

	d.topic = topic
	return nil
}

func (d *DeviceRegistry) CleanUp() error {
	if registry := d.GetRegistry(); registry != nil {
		_, err := d.client.Projects.Locations.Registries.Delete(d.RegistryName()).Do()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceRegistry) RegistryName() string {
	return fmt.Sprintf("projects/%s/locations/%s/registries/%s", d.projectID, d.Region, d.RegistryID)
}

func (d *DeviceRegistry) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Region: %s \n", d.Region))
	builder.WriteString(fmt.Sprintf("RegistryID: %s \n", d.RegistryID))
	builder.WriteString(fmt.Sprintf("RegistryName: %s \n", d.registry.Name))
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
		parent: DeviceRegistryParent(projectID, region),
	}
}