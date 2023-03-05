IMAGE_NAME ?= sonntuet1997/avalanche

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
	docker build --build-arg MODULE=worker . -t $(IMAGE_NAME):lastest

run-200-worker-docker:
	for i in {1..200}
	do
		docker run -d --name container-$i $(IMAGE_NAME)
	done

stop-200-worker-docker:
	for i in {1..200}
	do
		docker stop -d --name container-$i
	done