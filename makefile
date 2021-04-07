.PHONY: build
build:
	./hacks/build.sh

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: lint
lint:
	golangci-lint run --tests=false
	golangci-lint run --disable-all -E golint,goimports,misspell

.PHONY: run
run:
	cd ./cmd/shell-exporter/ && go run .