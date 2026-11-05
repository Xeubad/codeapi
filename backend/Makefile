.PHONY: build run generate test test-unit test-integration test-e2e

VERSION ?= $(shell ./scripts/resolve-version.sh)
LDFLAGS ?= -s -w -X main.Version=$(VERSION)

build:
	CGO_ENABLED=0 go build -tags embed -ldflags="$(LDFLAGS)" -trimpath -o bin/server ./cmd/server

# 本地起服务的唯一入口：先加载 .env（DATA_DIR 等），再启动 bin/server
run:
	@test -f .env || { echo "backend/.env 不存在：按 DEV_GUIDE「起服务姿势」创建后再运行"; exit 1; }
	set -a && . ./.env && set +a && exec ./bin/server

generate:
	go generate ./ent
	go generate ./cmd/server

test:
	go test ./...
	golangci-lint run ./...

test-unit:
	go test -tags=unit ./...

test-integration:
	go test -tags=integration ./...

test-e2e:
	./scripts/e2e-test.sh

test-e2e-local:
	go test -tags=e2e -v -timeout=300s ./internal/integration/...
