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

release:
	@echo "don't forget to create and push a git tag! e.g."
	@echo "  git tag -a v0.1.0 -m 'First release'"
	@echo "  git push origin v0.1.0"
	@sleep 3
	goreleaser release
.PHONY: release
