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

.PHONY: twitterwipe tweetstosql likeswipe likestosql all

all: twitterwipe tweetstosql likeswipe likestosql

twitterwipe:
	@echo ">  Building twitterwipe binary..." 
	$(GOBUILD) -v ./cmd/$(TWITTERWIPE)
	@echo ">  Done..." 
tweetstosql: 
	@echo ">  Building tweetstosql binary..." 
	$(GOBUILD) -v ./cmd/$(TWEETSTOSQL)
	@echo ">  Done..." 
likestosql:
	@echo ">  Building likestosql binary..." 
	$(GOBUILD) -v ./cmd/$(LIKESTOSQL)
	@echo ">  Done..." 
likeswipe:
	@echo ">  Building likeswipe binary..." 
	$(GOBUILD) -v ./cmd/$(LIKESWIPE)
	@echo ">  Done..." 

clean: 
	$(GOCLEAN)
	rm -f $(TWITTERWIPE)
	rm -f $(TWEETSTOSQL)
	rm -f $(LIKESTOSQL)
	rm -f $(LIKESWIPE)
