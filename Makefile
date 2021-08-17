dep:
	go mod tidy
	go mod vendor
test:
	go test -v ./...
build:
	go build
run: build
	go run main.go
docker-build:
	docker build .
start:
	docker-compose up -d --build
down:
	docker-compose down
