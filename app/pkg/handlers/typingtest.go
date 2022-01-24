package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type TypingTestRequest struct {
	Words int64 `json:"words"`
}

type TypingTestResponse struct {
	ErrorMsg string   `json:"error,omitempty"`
	Words []string `json:"word"`
}

type MonkeyTypeResponse struct {
	Words []string `json:"words"`
}

var wordsURL string

func init() {
	wordsURL = os.Getenv("TEST_WORDS_URL")
}

// TypingTest provides a typing test for discord.
func TypingTest(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var testRequest TypingTestRequest
	err := json.Unmarshal([]byte(request.Body), &testRequest)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to unmarshal: %v", err)}, err)
	}

	if testRequest.Words < 1 {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to request words less than 1.")}, nil)
	}

	if testRequest.Words > 20 {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprint("Invalid request words must be less then 20.")}, nil)
	}

	req, err := http.NewRequest("GET", wordsURL, nil)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to create request: %v", err)}, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to do request: %v", err)}, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body.")
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to read all body: %v", err)}, err)
	}

	var testResponse MonkeyTypeResponse
	err = json.Unmarshal(body, &testResponse)
	if err != nil {
		return apiResponse(http.StatusBadRequest, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to unmarshal data: %v", err)}, err)
	}

	lenWords := len(testResponse.Words)
	var typingTestRes TypingTestResponse
	for i := 0; i < int(testRequest.Words); i++ {
		number, err := rand.Int(rand.Reader, big.NewInt(int64(lenWords)))
		if err != nil {
			return apiResponse(http.StatusInternalServerError, TypingTestResponse{ErrorMsg: fmt.Sprintf("Failed to generate int: %v", err)}, err)
		}
		typingTestRes.Words = append(typingTestRes.Words, testResponse.Words[number.Int64()])
	}
	return apiResponse(http.StatusOK, typingTestRes, nil)
}
