package typechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type response struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []*choice `json:"choices"`
	Usage   *usage    `json:"usage"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type choice struct {
	Index        int      `json:"index"`
	Message      *message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type LanguageModel interface {
	complete(prompt string) (string, error)
}

func NewLanguageModel() (LanguageModel, error) {
	if os.Getenv("OPENAI_API_KEY") != "" {
		apiKey := os.Getenv("OPENAI_API_KEY")

		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			return nil, missingEnvironmentVariable("OPENAI_MODEL")
		}

		endPoint := os.Getenv("OPENAI_ENDPOINT")
		if endPoint == "" {
			endPoint = "https://api.openai.com/v1/chat/completions"
		}

		return NewOpenAILanguageModel(apiKey, model, endPoint, os.Getenv("OPENAI_ORGANIZATION")), nil
	}

	if os.Getenv("AZURE_OPENAI_API_KEY") != "" {
		apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
		endPoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
		if endPoint == "" {
			return nil, missingEnvironmentVariable("AZURE_OPENAI_ENDPOINT")
		}

		return NewAzureOpenAILanguageModel(apiKey, endPoint), nil
	}

	return nil, missingEnvironmentVariable("OPENAI_API_KEY or AZURE_OPENAI_API_KEY")
}

type Option func(*baseLanguageModel)

func WithHeaders(headers map[string]string) func(*baseLanguageModel) {
	return func(m *baseLanguageModel) {
		m.headers = headers
	}
}

func WithDefaultParams(defaultParams map[string]any) func(*baseLanguageModel) {
	return func(m *baseLanguageModel) {
		m.defaultParams = defaultParams
	}
}

type baseLanguageModel struct {
	url                string
	retryMaxAttempts   int
	retryPauseDuration time.Duration
	headers            map[string]string
	defaultParams      map[string]any
}

func newBaseLanguageModel(url string, options ...Option) LanguageModel {
	m := &baseLanguageModel{
		url:                url,
		retryMaxAttempts:   3,
		retryPauseDuration: 1000 * time.Millisecond,
		headers:            make(map[string]string),
	}

	for _, option := range options {
		option(m)
	}

	return m
}

func (m *baseLanguageModel) complete(prompt string) (string, error) {
	retryCount := 0

	for {
		paramMap := map[string]any{
			"messages": []map[string]any{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"temperature": 0,
			"n":           1,
		}

		if m.defaultParams != nil {
			for k, v := range m.defaultParams {
				paramMap[k] = v
			}
		}

		params, err := json.Marshal(paramMap)

		if err != nil {
			return "", err
		}

		req, err := http.NewRequest(http.MethodPost, m.url, bytes.NewBuffer(params))
		if err != nil {
			return "", err
		}

		req.Header.Add("Content-Type", "application/json;charset=utf-8")

		if m.headers != nil {
			for k, v := range m.headers {
				req.Header.Add(k, v)
			}
		}

		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		bodyBytes, err := m.readBody(resp)
		if err != nil {
			return "", err
		}

		if resp.StatusCode == http.StatusOK {

			var r response
			err = json.Unmarshal(bodyBytes, &r)
			if err != nil {
				return "", err
			}

			if len(r.Choices) == 0 || r.Choices[0].Message == nil {
				return "", nil
			}

			return r.Choices[0].Message.Content, nil
		}

		if !m.isTransientHttpError(resp.StatusCode) || retryCount >= m.retryMaxAttempts {
			return "", errors.New(fmt.Sprintf("REST API error %d: %s", resp.StatusCode, resp.Status))
		}

		time.Sleep(m.retryPauseDuration)
		retryCount++
	}
}

func (m *baseLanguageModel) readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

// 429: TooManyRequests
// 500: InternalServerError
// 502: BadGateway
// 503: ServiceUnavailable
// 504: GatewayTimeout
func (m *baseLanguageModel) isTransientHttpError(code int) bool {
	return code == 429 || code == 500 || code == 502 || code == 503 || code == 504
}

// NewAzureOpenAILanguageModel Creates a language model encapsulation of an Azure OpenAI REST API endpoint.
func NewAzureOpenAILanguageModel(apiKey, endPoint string) LanguageModel {
	return newBaseLanguageModel(endPoint, WithHeaders(map[string]string{
		"api-key": apiKey,
	}))
}

// NewOpenAILanguageModel Creates a language model encapsulation of an OpenAI REST API endpoint.
func NewOpenAILanguageModel(apiKey, model, endPoint, org string) LanguageModel {
	return newBaseLanguageModel(
		endPoint,
		WithHeaders(map[string]string{
			"Authorization":       fmt.Sprintf("Bearer %s", apiKey),
			"OpenAI-Organization": org,
		}),
		WithDefaultParams(map[string]any{
			"model": model,
		}),
	)
}

func missingEnvironmentVariable(name string) error {
	return fmt.Errorf("Missing environment variable: %s", name)
}
