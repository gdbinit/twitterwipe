-include .env

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TWITTERWIPE=twitterwipe
TWEETSTOSQL=tweetstosql
LIKESWIPE=likeswipe
LIKESTOSQL=likestosql
GOPATH := $(shell pwd)

.PHONY: twitterwipe tweetstosql likeswipe likestosql deps all

all: twitterwipe tweetstosql likeswipe likestosql

deps:
	@echo ">  Retrieving external packages..."
	@echo ">  github.com/dghubble/oauth1"
	@GOPATH=$(GOPATH) $(GOGET) github.com/dghubble/oauth1
	@echo ">  github.com/dghubble/go-twitter/twitter"
	@GOPATH=$(GOPATH) $(GOGET) github.com/dghubble/go-twitter/twitter
	@echo ">  github.com/lib/pq"
	@GOPATH=$(GOPATH) $(GOGET) github.com/lib/pq

twitterwipe:
	@echo ">  Building twitterwipe binary..." 
	@GOPATH=$(GOPATH) $(GOBUILD) -v cmd/$(TWITTERWIPE)
	@echo ">  Done..." 
tweetstosql: 
	@echo ">  Building tweetstosql binary..." 
	@GOPATH=$(GOPATH) $(GOBUILD) -v cmd/$(TWEETSTOSQL)
	@echo ">  Done..." 
likestosql:
	@echo ">  Building likestosql binary..." 
	@GOPATH=$(GOPATH) $(GOBUILD) -v cmd/$(LIKESTOSQL)
	@echo ">  Done..." 
likeswipe:
	@echo ">  Building likeswipe binary..." 
	@GOPATH=$(GOPATH) $(GOBUILD) -v cmd/$(LIKESWIPE)
	@echo ">  Done..." 

clean: 
	@GOPATH=$(GOPATH) $(GOCLEAN)
	rm -f $(TWITTERWIPE)
	rm -f $(TWEETSTOSQL)
	rm -f $(LIKESTOSQL)
	rm -f $(LIKESWIPE)
