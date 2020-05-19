# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

SERVICE_NAME = $(shell pwd |sed 's/^\(.*\)[/]//' )
PROJECT_NAME=${SERVICE_NAME}
BIN_DIR=/usr/sbin
CP_CMD=/usr/bin/cp
COMMAND=${SERVICE_NAME}

export GO111MODULE=on

all: build

build:
	go build -o $(PROJECT_NAME) -v  -ldflags "-X main.Version=$(version) -X main.GitCommit=`git rev-parse HEAD`"
test:
	$(GOTEST) -v ./...
install:build
	$(CP_CMD) $(COMMAND) $(DESTDIR)$(BIN_DIR)
clean:
	$(GOCLEAN)
	rm -f $(PROJECT_NAME)



