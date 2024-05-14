build: build/bitbucket-cli

build/bitbucket-cli:
	@mkdir -p build/
	@go build -o ./build/bitbucket-cli ./cmd/bitbucket-cli

clean:
	@rm -fr build/

fmt:
	@go fmt ./...
.PHONY: fmt
