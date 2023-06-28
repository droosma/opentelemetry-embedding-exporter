package embeddingexporter

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIEmbeddings struct {
	client *openai.Client
}

func NewOpenAIEmbeddings(key string, baseUri string, modelMapping map[string]string, version string) *OpenAIEmbeddings {
	if version == "" {
		version = "2023-05-15"
	}

	config := openai.DefaultAzureConfig(key, baseUri)
	config.APIVersion = version
	config.AzureModelMapperFunc = func(model string) string {
		return modelMapping[model]
	}
	return &OpenAIEmbeddings{client: openai.NewClientWithConfig(config)}
}

func (o *OpenAIEmbeddings) Generate(input string) (Embedding, error) {
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
