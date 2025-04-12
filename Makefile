default: build
run:
	@go run cmd/elements/main.go

clean:
	@rm -rf bin

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down

lint:
	golangci-lint run

url:
	@go run cmd/url/main.go
	
token:
	@go run cmd/token/main.go $(CODE)

fill:
	@go run cmd/fill/main.go