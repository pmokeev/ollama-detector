package user

import "context"

// RandomUserGenerator represents base interface for generating
// random user data.
type RandomUserGenerator interface {
	Generate(ctx context.Context, count int) []RandomUser
}
