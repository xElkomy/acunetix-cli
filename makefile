BINARY_NAME = acunetix-cli

build:
	GOARCH=amd64 GOOS=darwin go build -o ./target/$(BINARY_NAME)-darwin .
	GOARCH=amd64 GOOS=linux go build -o ./target/$(BINARY_NAME)-linux .
	GOARCH=amd64 GOOS=windows go build -o ./target/$(BINARY_NAME)-windows .

run: build
	./$(BINARY_NAME)

clean:
	go clean
	rm ./target/$(BINARY_NAME)-darwin
	rm ./target/$(BINARY_NAME)-linux
	rm ./target/$(BINARY_NAME)-windows
