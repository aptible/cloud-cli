GOPRIVATE?="github.com/aptible/*"

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
	docker run \
		--rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:v1.46 \
		golangci-lint run
.PHONY: lint

test: lint
	go test ./...
.PHONY: test

pretty:
	go fmt ./...
.PHONY: pretty
