# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=arma-admin
DB_PATH=$(HOME)/.arma-admin/arma_admin_db.sqlite

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
	./$(BINARY_NAME) makemigrations
	./$(BINARY_NAME) run
run-debug:
	rm -f $(DB_PATH)
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) makemigrations -loglevel=debug
	./$(BINARY_NAME) createsuperuser -loglevel=debug
	./$(BINARY_NAME) run -loglevel=debug


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v