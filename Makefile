.PHONY: build
build: fmt
	go build -o opap *.go

.PHONY: run
run: build
	./opap

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test
