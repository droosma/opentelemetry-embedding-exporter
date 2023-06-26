package embeddingexporter

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type embedding interface {
	Embed(input string) ([]float32, error)
}

type OpenAiEmbedder struct {
	client *openai.Client
}

func NewOpenAiEmbedder(key string,
	baseUri string,
	version ...string) *OpenAiEmbedder {
	v := "2023-05-15"
	if len(version) > 0 {
		v = version[0]
	}

	config := openai.DefaultAzureConfig(key, baseUri)
	config.APIVersion = v

	// config.AzureModelMapperFunc = func(model string) string {
	//    azureModelMapping = map[string]string{
	//        "gpt-3.5-turbo":"your gpt-3.5-turbo deployment name",
	//    }
	//    return azureModelMapping[model]
	// }
	return &OpenAiEmbedder{client: openai.NewClientWithConfig(config)}
}

func (o *OpenAiEmbedder) Embed(input string) ([]float32, error) {
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
