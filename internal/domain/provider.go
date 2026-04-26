package domain




type ProvideName string


const (
	ProviderUnknown ProvideName = "unknown"
	ProviderOpenAI  ProvideName = "openai"
	ProviderAnthropic ProvideName = "anthropic"
	ProviderMock ProvideName = "mock"
	ProviderOllama ProvideName = "ollama"
)

type GenerateRequest struct {
	RequestID string 
	Model string
	Messages []Message
	ResponseFormat string
	Metadata map[string]string
	
}

type GenerateResponse struct{
	Provider ProvideName
	Model string 
	Output map[string]any
}