# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=todoList

all: 
	go mod tidy
	make format
	make clean 
	make build
	./$(BINARY_NAME)

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy:
	docker build -t http .
	docker run -p 8080:8080 http

format:
	$(GOFMT) ./...