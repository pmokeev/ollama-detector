package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog"
)

const (
	// _randomuserBaseURL represents base URL of randomuser.me API service
	// for generating random customers data.
	_randomuserBaseURL = "https://randomuser.me/api/"
)

// UserGenerator represents struct for generating random user data.
type UserGenerator struct {
	logger *zerolog.Logger
}

// NewUserGenerator returns a new instance of UserGenerator.
func NewUserGenerator(logger *zerolog.Logger) *UserGenerator {
	log := logger.With().
		Str("component", "user-generator").
		Logger()

	return &UserGenerator{
		logger: &log,
	}
}

// Generate generates user data in the amount passed by the parameter.
func (c *UserGenerator) Generate(ctx context.Context, count int) []RandomUser {
	randomUserURL, err := url.Parse(_randomuserBaseURL)
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while parsing random user base URL")
	}

	randomUserURL.RawQuery = url.Values{
		"results": {strconv.Itoa(count)},
	}.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, randomUserURL.String(), nil)
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while getting random users")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while doing http request")
	}
	defer response.Body.Close()

	var randomUserResponse RandomUserResponse
	if err := json.NewDecoder(response.Body).Decode(&randomUserResponse); err != nil {
		c.logger.Fatal().
			Err(err).
			Msg("Error occured while decoding http response")
	}

	return randomUserResponse.Results
}
