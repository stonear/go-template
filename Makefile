BINARY_NAME=go-template

build:
	gen
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

vet:
	go vet

lint:
	golangci-lint run
	go mod tidy
	@if ! git diff --quiet; then \
		echo "'go mod tidy' resulted in changes or working tree is dirty:"; \
		git --no-pager diff; \
	fi

gen:
	cd db; sqlc generate;
	swag init --parseDependency --parseInternal --parseDepth 1
