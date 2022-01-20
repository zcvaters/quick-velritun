package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var wordsUrl = "https://monkeytype.com/languages/english.json"

type TypingTestRequest struct {
	Words int `json:"words"`
}

type TypingTestResponse struct {
	Error string   `json:"error,omitempty"`
	Words []string `json:"word"`
}

type MonkeyTypeResponse struct {
	Words []string `json:"words"`
}

// TypingTest for creating a request to create a typing test.
func TypingTest(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var testRequest TypingTestRequest
	err := json.Unmarshal([]byte(request.Body), &testRequest)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprintf("Failed to unmarshal: %v", err)})
	}

	if testRequest.Words > 20 {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprint("Invalid request words must be less then 20.")})
	}

	req, err := http.NewRequest("GET", wordsUrl, nil)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprintf("Failed to create request: %v", err)})
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprintf("Failed to do request: %v", err)})
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprintf("Failed to read all body: %v", err)})
	}

	var testResponse MonkeyTypeResponse
	err = json.Unmarshal(body, &testResponse)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{Error: fmt.Sprintf("Failed to unmarshal data: %v", err)})
	}

	lenWords := len(testResponse.Words)
	var typingTestRes TypingTestResponse
	for i := 0; i < testRequest.Words; i++ {
		number, err := rand.Int(rand.Reader, big.NewInt(int64(lenWords)))
		if err != nil {
			return apiResponse(http.StatusInternalServerError, TypingTestResponse{Error: fmt.Sprintf("Failed to generate int: %v", err)})
		}
		typingTestRes.Words = append(typingTestRes.Words, testResponse.Words[number.Int64()])
	}
	return apiResponse(http.StatusOK, typingTestRes)
}

// func main() {
// 	lambda.Start(handler)
// }
