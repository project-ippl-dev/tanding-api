package config

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type StorageClient struct {
	Client *storage.Client
	Config StorageConfig
}

func NewStorageClient(conf StorageConfig) StorageClient {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if err != nil {
		log.Fatalf("%v", err)
	}

	storageClient, err := storage.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("%v", err)
	}

	return StorageClient{
		Client: storageClient,
		Config: conf,
	}
}
