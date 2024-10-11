package ollama

import "context"

// LLMChatter represents base interface for LLM interaction.
type LLMChatter interface {
	Chat(ctx context.Context, text string) string
}
