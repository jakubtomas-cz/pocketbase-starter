RUN_COMMAND := go run cmd/server/main.go serve --http=localhost:8090

.PHONY: run build dev docker

run:
	${RUN_COMMAND}

build:
	go build -o bin/app ./cmd/server/main.go

dev:
	npx nodemon --signal SIGTERM --watch . --ext go --ignore pb_data/ --ignore pb_public/ --ignore bin/ --exec ${RUN_COMMAND}

docker:
	docker build -t pocketbase-starter:${IMAGE_TAG} .
