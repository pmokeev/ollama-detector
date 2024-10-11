package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

const (
	// _ollamaModelName represents name of ollama model.
	_ollamaModelName = "llama3"
	// _ollamaRole represents ollama role for quering to Ollama.
	_ollamaRole = "user"
	// _ollamaBaseURL represents base URL for connecting to Ollama.
	_ollamaBaseURL = "http://localhost:11434/api/chat"
)

// Client represents client for chatting with Ollama LLM.
type Client struct {
	logger *zerolog.Logger
}

// NewClient returns a new instance of Ollama client.
func NewClient(logger *zerolog.Logger) *Client {
	log := logger.With().
		Str("component", "ollama-client").
		Logger()

	return &Client{
		logger: &log,
	}
}

// Chat sends text message to Ollama LLM and returns generated message.
func (c *Client) Chat(ctx context.Context, text string) string {
	ollamaRequest := OllamaRequest{
		Model: _ollamaModelName,
		Messages: []Message{
			{
				Role:    _ollamaRole,
				Content: text,
			},
		},
	}

	body, err := json.Marshal(&ollamaRequest)
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while marhalling request")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, _ollamaBaseURL, bytes.NewReader(body))
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while doing http request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while doing http request")
	}
	defer response.Body.Close()

	var ollamaResponse OllamaResponse
	if err := json.NewDecoder(response.Body).Decode(&ollamaResponse); err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while decoding http request")
	}

	return ollamaResponse.Message.Content
}
