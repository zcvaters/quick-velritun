package tests

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	handlers "github.com/zcvaters/quick-velritun/app/pkg/handlers"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {
		res, err := handlers.TypingTest(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
	})
}
