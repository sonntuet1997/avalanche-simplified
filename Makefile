up:
	docker-compose up

down:
	docker-compose down

lint:
	golint -set_exit_status ./...

love:
	- cd src/worker && gosec ./...

tidy:
	cd src/worker && go mod tidy

test:
	cd src/worker && go mod tidy && go test ./...

run-worker:
	cd src/worker && go run .

build-worker:
	cd src/worker && go run .

build-worker-docker:
	cd src/worker && go run .
