VERSION ?= $(shell git describe --tags --dirty --match 'v*.*.*')

.PHONY: build
build:
	mkdir -p dist
	go build \
		-v \
		-race \
		-mod vendor \
		-o dist \
		-ldflags "-X main.version=$(VERSION)" \
		./cmd/...

.PHONY: install
install:
	go install \
		-v \
		-race \
		-mod vendor \
		-ldflags "-X main.version=$(VERSION)" \
		./cmd/...

.PHONY: clean
clean:
	rm -rf dist
	rm -rf cover
	rm -rf prof
	rm -rf passgen.test

.PHONY: test
test:
	mkdir -p cover
	go test \
		-v \
		-race \
		-coverprofile cover/cover.out \
		./...

.PHONY: cover
cover:
	go tool cover \
		-html \
		cover/cover.out

.PHONY: vet
vet:
	go vet -v ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: benchmark
benchmark:
	mkdir -p prof
	go test \
		-v \
		-run Benchmark \
		-cpuprofile prof/cpu.prof \
		-memprofile prof/mem.prof \
		-bench ./...
