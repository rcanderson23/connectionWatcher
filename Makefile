BIN ?= connectionWatcher
IMG ?= ghcr.io/rcanderson23/connectionwatcher:latest

.PHONY: test
test: fmt vet
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: dep
dep:
	go mod download

.PHONY: build
build: dep fmt vet
	go build -ldflags='-w -extldflags "-static"' -o bin/$(BIN) main.go

.PHONY: docker-push
docker-build:
	docker build -t=$(IMG) .

.PHONY: docker-push
docker-push:
	docker push $(IMG)
