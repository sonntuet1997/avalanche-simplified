IMAGE_NAME ?= sonntuet1997/avalanche

lint:
	golint -set_exit_status ./...

love:
	- cd src/worker && gosec ./...

tidy:
	cd src/worker && go mod tidy

test:
	cd src/worker && go mod tidy && go test ./...

test-e2e:
	cd src/worker/functional-tests && go mod tidy && OOS=linux GOARCH=amd64 CGO_ENABLED=0 go test ./... -v

vet:
	cd src/worker && go vet ./... && staticcheck ./...

run-worker:
	cd src/worker && go run .

build-docker:
	docker-compose build

run-200-worker-docker:
	docker-compose up -d --scale node=200

down-200-worker-docker:
	docker-compose down