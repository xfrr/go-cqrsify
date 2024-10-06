# This makefile is used for development purposes only.

.PHONY: bench
bench:
	go test -bench=. -benchmem -count=1 -run=^$$ ./...

.PHONY: cover-html
cover-html: cover
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test -v ./...

.PHONY: cover-out
test-cover-out: cover
	go tool cover -func=coverage.out

.PHONY: cover
test-cover: 
	go test -coverprofile=coverage.out ./...
