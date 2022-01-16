package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil && err.Error() != ErrNon200Response.Error() {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Unable decode IP", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid HTTP response")
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})

	t.Run("Successful Request", func(t *testing.T) {

		typeRequest := TypingTestRequest{Words: 10}
		typeReqJson, _ := json.Marshal(typeRequest)

		req := events.APIGatewayProxyRequest{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    string(typeReqJson),
		}
		res, err := handler(req)
		if err != nil {
			t.Fatal("failed to make typing test request")
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("failed typing test request statusCode %v, expected %v", res.StatusCode, http.StatusOK)
		}
		var resWords TypingTestResponse
		err = json.Unmarshal([]byte(res.Body), &resWords)
		if err != nil {
			t.Fatal("failed to unmarshal json")
		}
		fmt.Println(resWords.Words)

	})
}
