.PHONY: format
format:
	gofmt -w -s .
	goimports -l -w .
	gofumpt -l -w .
	gci write .

.PHONY: lint
lint:
	@golangci-lint run --sort-results

.PHONY: test
test:
	go test -v -p 1 ./...

build_server:
	docker build --cache-from yvv4docker/task-wow-server:latest -t yvv4docker/task-wow-server:latest -f server.Dockerfile .

build_client:
	docker build --cache-from yvv4docker/task-wow-client:latest -t yvv4docker/task-wow-client:latest -f client.Dockerfile .

docker_push:
	docker tag yvv4docker/task-wow-server:latest yvv4docker/task-wow-server:v0.0.1
	docker push yvv4docker/task-wow-server:v0.0.1

up:
	docker-compose up

down:
	docker-compose down
