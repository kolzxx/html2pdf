CMD_PKG=./cmd/main
BINARY_NAME=server

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

all: build
.PHONY: all

build: clean swagger
	go build -ldflags="-w -linkmode=external -X main.Version=$(VERSION)" -v -o $(BINARY_NAME) $(CMD_PKG)
.PHONY: build

clean:
	rm -f ./$(BINARY_NAME)
.PHONY: clean

tidy:
	go mod tidy
.PHONY: tidy

# @see https://github.com/swaggo/swag/
swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.8.10
	swag init --parseDependency --parseInternal -d cmd/main,internal/controllers -o api --ot go,json --md api
.PHONY: swagger
