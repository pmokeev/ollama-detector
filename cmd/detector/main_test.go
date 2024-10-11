package main

import (
	"context"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/pmokeev/ollama-detector/pkg/ollama"
	"github.com/pmokeev/ollama-detector/pkg/user"
	"github.com/stretchr/testify/assert"
)

func Test_getCompanyInfo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Company
	}{
		{
			name:  "Success: single product",
			input: "Google\nApples\nUSA\n",
			want: Company{
				Name:     "Google",
				Products: []string{"Apples"},
				Market:   "USA",
			},
		},
		{
			name:  "Success: multiple products",
			input: "Google\nApples,Bananas\nUSA\n",
			want: Company{
				Name:     "Google",
				Products: []string{"Apples", "Bananas"},
				Market:   "USA",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinDefer, err := mockStdin(t, tt.input)
			if err != nil {
				t.Fatal(err)
			}

			defer stdinDefer()

			got := getCompanyInfo()

			assert.True(t, reflect.DeepEqual(got, tt.want))
		})
	}
}

func Test_fetchCustomerData(t *testing.T) {
	type args struct {
		ctx   context.Context
		count int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success: single customer",
			args: args{
				ctx:   context.TODO(),
				count: 1,
			},
		},
		{
			name: "Success: multiple customers",
			args: args{
				ctx:   context.TODO(),
				count: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_userGenerator = user.NewUserGenerator(&_logger)

			got := fetchCustomerData(tt.args.ctx, tt.args.count)

			assert.Equal(t, tt.args.count, len(got))
			for _, customer := range got {
				assert.NotEmpty(t, customer.Name)
				assert.NotEmpty(t, customer.Age)
				assert.NotEmpty(t, customer.Email)
			}
		})
	}
}

func Test_analyzeOpportunities(t *testing.T) {
	type args struct {
		ctx       context.Context
		company   Company
		customers []Customer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success: singular customer",
			args: args{
				ctx: context.TODO(),
				company: Company{
					Name:     "Google",
					Products: []string{"Apples", "Bananas"},
					Market:   "USA",
				},
				customers: []Customer{
					{
						Name:  "Google Google",
						Age:   30,
						Email: "google@google.com",
					},
				},
			},
		},
		{
			name: "Success: multiple customers",
			args: args{
				ctx: context.TODO(),
				company: Company{
					Name:     "Google",
					Products: []string{"Apples", "Bananas"},
					Market:   "USA",
				},
				customers: []Customer{
					{
						Name:  "Google Google",
						Age:   30,
						Email: "google@google.com",
					},
					{
						Name:  "Yahoo Yahoo",
						Age:   20,
						Email: "yahoo@yahoo.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_llmChatter = ollama.NewClient(&_logger)

			got := analyzeOpportunities(tt.args.ctx, tt.args.company, tt.args.customers)

			assert.Equal(t, len(got), len(tt.args.customers))

			ind := 0
			for customer, opportunity := range got {
				assert.Equal(t, customer, tt.args.customers[ind].String())
				assert.NotEmpty(t, opportunity)
				ind += 1
			}
		})
	}
}

func Test_outputResults(t *testing.T) {
	type args struct {
		outputType    OutputType
		opportunities map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				outputType: JSONOutputType,
				opportunities: map[string]string{
					"Name": "Opportunity",
				},
			},
			want: "{\"Name\":\"Opportunity\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			os.Stdout = w

			outputResults(tt.args.outputType, tt.args.opportunities)

			w.Close()
			out, _ := io.ReadAll(r)

			assert.Equal(t, tt.want, string(out))
		})
	}
}

func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin

	tmpfile, err := os.CreateTemp("", "temp")
	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	os.Stdin = tmpfile

	return func() {
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}
