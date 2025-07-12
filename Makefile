.PHONY: build
build:
	go build -v .
	gh extension remove issue-clone
	gh extension install .
