# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=todoList

all: 
	make format
	make clean 
	go run main.go
	
build:
	make format
	make clean 
	$(GOBUILD) -o $(BINARY_NAME) -v
	echo "preview on http://localhost:8080"
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy:
	docker build -t http .
	docker run -p 8080:8080 http

format:
	$(GOFMT) ./...