build: build/bitbucket-cli

build/bitbucket-cli:
	mkdir -p build/
	go build -o ./build/bitbucket-cli ./cmd/bitbucket-cli

clean:
	rm -rf build/
