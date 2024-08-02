.PHONY: lint test test-all test-cov

lint:
	@golangci-lint run \
		--enable-all \
		--disable execinquery \
		--disable gomnd \
		--disable tagliatelle

test:
	@go test -race -short ./...

test-all:
	@go test -race ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out
