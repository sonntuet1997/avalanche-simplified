IMAGE_NAME ?= sonntuet1997/avalanche

lint:
	golint -set_exit_status ./...

love:
	- cd src/worker && gosec ./...

tidy:
	cd src/worker && go mod tidy

test:
	cd src/worker && go mod tidy && go test ./...

vet:
	cd src/worker && go vet ./... && staticcheck ./...

run-worker:
	cd src/worker && go run .

build-worker:
	docker-compose build
	docker build --build-arg MODULE=worker . -t $(IMAGE_NAME):lastest

run-200-worker-docker:


stop-200-worker-docker:
