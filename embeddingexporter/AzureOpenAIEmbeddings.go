package embeddingexporter

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type AzureOpenAIEmbeddings struct {
	client *openai.Client
}

func NewAzureOpenAIEmbeddings(key string, baseUri string, modelMapping map[string]string, version string) *AzureOpenAIEmbeddings {
	if version == "" {
		version = "2023-05-15"
	}

	config := openai.DefaultAzureConfig(key, baseUri)
	config.APIVersion = version
	config.AzureModelMapperFunc = func(model string) string {
		return modelMapping[model]
	}
	return &AzureOpenAIEmbeddings{client: openai.NewClientWithConfig(config)}
}

func (o *AzureOpenAIEmbeddings) Generate(input string) (Embedding, error) {
	resp, err := o.client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{input},
			Model: openai.AdaEmbeddingV2,
		})

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}
