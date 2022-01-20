package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	handlers "github.com/zcvaters/quick-velritun/app/pkg/handlers"
)

func main() {
	lambda.Start(handlers.TypingTest)
}
