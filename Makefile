.PHONY: build

build:
	sam build
	sam local start-api
