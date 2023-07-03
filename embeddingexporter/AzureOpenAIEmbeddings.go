package embeddingexporter

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

type AzureOpenAIEmbeddings struct {
	client *azopenai.Client
}

func NewAzureOpenAIEmbeddings(key string, endpoint string, modelId string, version string) *AzureOpenAIEmbeddings {
	cred, err := azopenai.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	options := &azopenai.ClientOptions{}
	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, modelId, options)

	if err != nil {
		panic(err)
	}

	return &AzureOpenAIEmbeddings{client: client}
}

func (o *AzureOpenAIEmbeddings) Generate(input string) (Embedding, error) {
	ctx := context.Background()
	body := azopenai.EmbeddingsOptions{
		Input: []string{input},
	}
	resp, err := o.client.GetEmbeddings(ctx, body, nil)

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}
