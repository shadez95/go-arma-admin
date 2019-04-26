# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=arma-admin

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./${BINARY_NAME} makemigrations
	./$(BINARY_NAME) run
run-debug:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./${BINARY_NAME} makemigrations
	./${BINARY_NAME} createsuperuser
	./$(BINARY_NAME) run -loglevel=debug


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v