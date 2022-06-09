GOPRIVATE?="github.com/aptible/cloud-api-clients"

init:
	GOPRIVATE=$(GOPRIVATE) go mod download
.PHONY: init

build:
	go build -o build/aptible
.PHONY: build

start:
	go run main.go
.PHONY: start

stop:
	echo "nothing to stop"
.PHONY: stop

destroy:
	echo "nothing to destroy"
.PHONY: destroy

clean:
	rm build/aptible
.PHONY: clean

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -E goimports -E godot
.PHONY: lint

test: lint
	go test ./...
.PHONY: test

pretty:
	go fmt ./...
.PHONY: pretty
