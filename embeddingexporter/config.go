package embeddingexporter

type Config struct {
	OpenAiKey     string
	OpenAiUri     string
	OpenAiVersion string
}

func NewConfig() Config {
	return Config{
		OpenAiVersion: "2023-05-15",
		OpenAiKey:     "",
		OpenAiUri:     "https://rg-openai-sandbox.openai.azure.com/",
	}
}
