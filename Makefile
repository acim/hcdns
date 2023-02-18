.PHONY: lint test test-all test-cov

lint:
	@golangci-lint run \
		--enable-all \
		--disable deadcode \
		--disable exhaustivestruct \
		--disable golint \
		--disable ifshort \
		--disable interfacer \
		--disable maligned \
		--disable nosnakecase \
		--disable scopelint \
		--disable structcheck \
		--disable tagliatelle \
		--disable varcheck

test:
	@go test -race -short ./...

test-all:
	@go test -race ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out
