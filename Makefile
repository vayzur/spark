BINARY_NAME=spark
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GORUN=$(GO) run
GOFMT=$(GO) fmt
GOMOD=$(GO) mod

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-w -s" .

run:
	$(GORUN) cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	$(GOFMT) ./...

deps:
	$(GOMOD) download

install:
	$(GOMOD) tidy

lint:
	golint ./...

vet:
	$(GO) vet ./...

check: lint vet test

.PHONY: all build run test clean fmt deps install lint vet check dev start
