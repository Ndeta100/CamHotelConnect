build:
	@go build -o bin/api

run:build
	@./bin/api

test:
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f /bin/api