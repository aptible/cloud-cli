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
	go vet ./...
.PHONY: lint

test:
	go test ./...
.PHONY: test

pretty:
	go fmt ./...
.PHONY: pretty
