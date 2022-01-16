package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
)

var (

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 Response found")

	wordsUrl = "https://monkeytype.com/languages/english.json"
)

type TypingTestRequest struct {
	Words int `json:"words,omitempty"`
}

type TypingTestResponse struct {
	Words []string `json:"words"`
}

type MonkeyTypeResponse struct {
	Words []string `json:"words"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var testRequest TypingTestRequest
	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &testRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("failed to make request %v", err.Error()), StatusCode: http.StatusBadRequest}, nil
	}

	req, _ := http.NewRequest("GET", wordsUrl, nil)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	var testResponse MonkeyTypeResponse
	err = json.Unmarshal(body, &testResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

	lenWords := len(testResponse.Words)
	var typingTestRes TypingTestResponse
	for i := 0; i < testRequest.Words; i++ {
		number, err := rand.Int(rand.Reader, big.NewInt(int64(lenWords)))
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		typingTestRes.Words = append(typingTestRes.Words, testResponse.Words[number.Int64()])
	}

	stringBody, _ := json.Marshal(typingTestRes)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(stringBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
