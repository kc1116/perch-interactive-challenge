package core

import (
	"cloud.google.com/go/pubsub"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudiot/v1"
	"sync"
)

var (
	httpClientOnce   sync.Once
	pubSubClientOnce sync.Once
	gcClient         *cloudiot.Service
	pubSubClient     *pubsub.Client
)

// GCHttpClient initializes google cloud Client from credentials in env
func GCHttpClient() (*cloudiot.Service, error) {
	var clientErr error
	httpClientOnce.Do(func() {
		ctx := context.Background()
		httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
		if err != nil {
			clientErr = err
			return
		}

		gcClient, err = cloudiot.New(httpClient)
		if err != nil {
			clientErr = err
			return
		}

	})

	return gcClient, clientErr
}

// PubSubClient initializes google cloud pubsub Client from credentials in env
func PubSubClient(projectID string) (*pubsub.Client, error) {
	var clientErr error
	pubSubClientOnce.Do(func() {
		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			clientErr = err
			return
		}

		pubSubClient = client
		return
	})

	return pubSubClient, clientErr
}
