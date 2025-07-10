.PHONY: build
build:
	go build .
	gh extension remove issue-clone
	gh extension install .
