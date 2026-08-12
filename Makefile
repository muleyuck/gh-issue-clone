.PHONY: test
test:
	golangci-lint run ./...
	go test -v ./...
	go build -v .

.PHONY: build
build:
	go build -v .
	gh extension remove issue-clone
	gh extension install .
