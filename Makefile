.PHONY: build

include .env
export

build:
	sam build
	sam local start-api --parameter-overrides WordsURL=${TEST_WORDS_URL}

test:
	go test ./...
	