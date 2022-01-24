package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	handlers "github.com/zcvaters/quick-velritun/app/pkg/handlers"
)

func TestTypingTestHandler(t *testing.T) {
	type want struct {
		Error     string
		WordCount int
	}

	tests := []struct {
		Want  want
		Input handlers.TypingTestRequest
	}{
		{Input: handlers.TypingTestRequest{Words: -1}, Want: want{WordCount: 0, Error: "Failed to request words less than 1."}},
		{Input: handlers.TypingTestRequest{Words: 5}, Want: want{WordCount: 5, Error: ""}},
		{Input: handlers.TypingTestRequest{Words: 10}, Want: want{WordCount: 10, Error: ""}},
		{Input: handlers.TypingTestRequest{Words: 15}, Want: want{WordCount: 15, Error: ""}},
		{Input: handlers.TypingTestRequest{Words: 20}, Want: want{WordCount: 20, Error: ""}},
		{Input: handlers.TypingTestRequest{Words: 25}, Want: want{WordCount: 0, Error: "Invalid request words must be less then 20."}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Words:%v", tc.Input.Words), func(t *testing.T) {
			tcJSON, err := json.Marshal(tc.Input)
			if err != nil {
				t.Fatal(err)
			}
			req := events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Path:       "/typingtest",
				Body:       string(tcJSON),
			}
			res, err := handlers.TypingTest(req)
			if err != nil {
				t.Fatal(err)
			}
			var actual handlers.TypingTestResponse
			json.Unmarshal([]byte(res.Body), &actual)
			if len(actual.Words) != tc.Want.WordCount {
				t.Fatalf("Expected word count %v got: %v", tc.Want.WordCount, len(actual.Words))
			}
			if actual.ErrorMsg != tc.Want.Error {
				t.Fatalf("Expected error %v got: %v", tc.Want.Error, actual.ErrorMsg)
			}
		})
	}
}
