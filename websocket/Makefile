# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=notify

all: 
	go mod tidy
	make format
	go run main.go

build:
	go mod tidy
	make format
	make clean 
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy:
	docker build -t $(BINARY_NAME) .
	docker run -p 1234:1234 $(BINARY_NAME)

format:
	$(GOFMT) ./...