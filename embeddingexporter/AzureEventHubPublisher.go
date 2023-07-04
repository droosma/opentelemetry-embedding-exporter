package embeddingexporter

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

type AzureEventHubPublisher struct {
	client *azeventhubs.ProducerClient
}

func NewAzureEventHubPublisher(connectionString string) *AzureEventHubPublisher {

	client, err := azeventhubs.NewProducerClientFromConnectionString(connectionString, "", nil)

	if err != nil {
		panic(err)
	}

	return &AzureEventHubPublisher{client: client}
}

func (a *AzureEventHubPublisher) Publish(embeddings []logEntryWithEmbedding) error {
	ctx := context.Background()
	newBatchOptions := &azeventhubs.EventDataBatchOptions{}

	batch, err := a.client.NewEventDataBatch(ctx, newBatchOptions)
	if err != nil {
		return err
	}

	batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("hello"),
	}, nil)

	return a.client.SendEventDataBatch(ctx, batch, nil)
}
