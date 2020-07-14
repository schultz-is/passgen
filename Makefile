PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: build
build:
	mkdir -p dist
	go build \
		-v \
		-race \
		-mod vendor \
		-o dist \
		./cmd/...

.PHONY: clean
clean:
	rm -rf dist

.PHONY: test
test:
	go test -v -race $(PKGS)

.PHONY: vet
vet:
	go vet -v $(PKGS)

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: staticcheck
staticcheck:
	staticcheck $(PKGS)
