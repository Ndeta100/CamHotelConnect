build:
	@go build -o hotel/api

run:build
	@./hotel/api

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go

drop:
	@go run scripts/cleanDB.go

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ./hotel

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 api