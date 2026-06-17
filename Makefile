.PHONY: build test dev deploy swag clean setup

BINARY = bootstrap

setup:
	go mod download
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example"; else echo ".env already exists, skipping"; fi

swag:
	swag init -g cmd/api/main.go -o docs/swagger --parseDependency

build: swag
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BINARY) ./cmd/api

test:
	go test ./... -v -count=1

dev:
	go run ./cmd/api

deploy: build
	zip function.zip $(BINARY)
	aws lambda update-function-code \
		--function-name agentic-serverless-api \
		--zip-file fileb://function.zip; \
	rm -f function.zip

clean:
	rm -f $(BINARY) function.zip
