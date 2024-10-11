package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pmokeev/ollama-detector/pkg/ollama"
	"github.com/pmokeev/ollama-detector/pkg/user"
	"github.com/rs/zerolog"
)

// OutputType represents type of output result.
type OutputType string

var (
	// _userGenerator represents contract for user generating.
	_userGenerator user.RandomUserGenerator
	// _llmChatter represents contract for LLM chatting.
	_llmChatter ollama.LLMChatter
	// _logger represents logger for printing warn/fatal/error/info logs.
	_logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	// JSONOutputType represents name of json output type.
	JSONOutputType OutputType = "json"
	// CSVOutputType represents name of csv output type.
	CSVOutputType OutputType = "csv"
)

// Company represents information about provided company.
type Company struct {
	Name     string
	Products []string
	Market   string
}

// String implements Stringer interface for company struct.
func (c *Company) String() string {
	return fmt.Sprintf(
		"Company name is %s, which produces products like %s and aims to %s market",
		c.Name,
		strings.Join(c.Products, ","),
		c.Market,
	)
}

// Customer represents information about generated customer.
type Customer struct {
	Name  string
	Age   int
	Email string
}

// String implements Stringer interface for customer struct.
func (c *Customer) String() string {
	return fmt.Sprintf(
		"Customer name is %s, his age is %d, and contact email: %s",
		c.Name,
		c.Age,
		c.Email,
	)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		s := <-signalChan
		_logger.Warn().
			Str("signal", s.String()).
			Msg("Got signal")
		cancel()
	}()

	_userGenerator = user.NewUserGenerator(&_logger)
	_llmChatter = ollama.NewClient(&_logger)

	company := getCompanyInfo()
	customers := fetchCustomerData(ctx, 1)
	opportunities := analyzeOpportunities(ctx, company, customers)
	outputResults(JSONOutputType, opportunities)
}

// getCompanyInfo gets company information from stdin.
func getCompanyInfo() Company {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please, enter the name of the company:\n")
	name, err := reader.ReadString('\n')
	if err != nil {
		_logger.Fatal().
			Err(err).
			Msg("Error occured while reading from stdin")
	}
	name = strings.ReplaceAll(name, "\n", "")

	fmt.Print("Please, enter products of the company (use `,` as separator):\n")
	productsInput, err := reader.ReadString('\n')
	if err != nil {
		_logger.Fatal().
			Err(err).
			Msg("Error occured while reading from stdin")
	}
	productsInput = strings.ReplaceAll(productsInput, "\n", "")
	produts := strings.Split(productsInput, ",")

	fmt.Print("Please, enter the name of the market the company is aims to:\n")
	market, err := reader.ReadString('\n')
	if err != nil {
		_logger.Fatal().
			Err(err).
			Msg("Error occured while reading from stdin")
	}
	market = strings.ReplaceAll(market, "\n", "")

	return Company{
		Name:     name,
		Products: produts,
		Market:   market,
	}
}

// fetchCustomerData fetches customer data using some random generator service.
func fetchCustomerData(ctx context.Context, count int) []Customer {
	randomUsers := _userGenerator.Generate(ctx, count)

	customers := make([]Customer, 0)
	for _, user := range randomUsers {
		customers = append(customers, Customer{
			Name:  fmt.Sprintf("%s %s", user.RandomUserName.First, user.RandomUserName.Last),
			Age:   user.RandomUserDob.Age,
			Email: user.Email,
		})
	}

	return customers
}

// analyzeOpportunities analyzes and returns opportunities using provided company and customers information.
func analyzeOpportunities(ctx context.Context, company Company, customers []Customer) map[string]string {
	opportunities := make(map[string]string, 0)

	for _, customer := range customers {
		text := fmt.Sprintf(
			"Use provided information about Company and generate selling suggestion for customer:\nCompany: %s\nCustomer: %s",
			company.String(),
			customer.String(),
		)

		opportunity := _llmChatter.Chat(ctx, text)
		opportunities[customer.String()] = opportunity
	}

	return opportunities
}

// outputResults prints structured opportunities.
func outputResults(outputType OutputType, opportunities map[string]string) {
	switch outputType {
	case JSONOutputType:
		output, err := json.Marshal(opportunities)
		if err != nil {
			_logger.Fatal().
				Err(err).
				Msg("Error occured while marhalling output")
		}
		fmt.Println(string(output))
	case CSVOutputType:
		// Actually it has to be rewritten to be more informative, because
		// now stdout contains canvas of text.
		for customer, opportunity := range opportunities {
			fmt.Printf("%s,%s", customer, opportunity)
		}
	default:
		_logger.Fatal().
			Str("outputType", string(outputType)).
			Msg("Provided unexpected outputType")
	}
}
