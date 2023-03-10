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

build-docker:
	docker-compose build

run-200-worker-docker:
	docker-compose up -d --scale node=20

stop-200-worker-docker:
	docker-compose stop