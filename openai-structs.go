package main

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	LogProbs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptToken      int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type promptRequestPayload struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type chatGPTResponseBody struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int64       `json:"created"`
	Model             string      `json:"model"`
	Choices           []Choice    `json:"choices"`
	Usage             Usage       `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}
