-include .env

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TWITTERWIPE=twitterwipe

.PHONY: twitterwipe all

all: twitterwipe

twitterwipe:
	@echo ">  Building twitterwipe binary..." 
	$(GOBUILD) -o $(TWITTERWIPE) main.go 
	@echo ">  Done..." 

clean: 
	$(GOCLEAN)
	rm -f $(TWITTERWIPE)
