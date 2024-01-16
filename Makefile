include .env
export
.SILENT:
WEBSITE_REPO=gitlab.com/a5805/ondeu/ondeu-back
CONTAINER_REGISTER=registry.gitlab.com/a5805/ondeu/ondeu-back
MOCKS_DESTINATION=internal/mocks
VERSION?="0.0.1"

app:
	go run cmd/api/main.go

install:
	go get -d -v ./.../

build:
	go build -o ./main ./cmd/api/main.go

test:
	go test ./.../ -v

docker_dev:
	docker buildx build \
	--progress=plain \
    --push \
    --cache-from type=local,src=./cache/back \
    --platform linux/amd64,linux/arm64 \
    --cache-to type=local,dest=./cache/back \
    -t ${CONTAINER_REGISTER}:dev .

docker_prod:
	docker buildx build \
	--progress=plain \
    --push \
    --cache-from type=local,src=./cache/back \
    --platform linux/amd64,linux/arm64 \
    --cache-to type=local,dest=./cache/back \
    -t  ${CONTAINER_REGISTER}:prod .

.NOTPARALLEL:

.PHONY: app mocks