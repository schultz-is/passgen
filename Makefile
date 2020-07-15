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
	rm -rf cover

.PHONY: test
test:
	mkdir -p cover
	go test \
		-v \
		-race \
		-coverprofile cover/cover.out \
		$(PKGS)

.PHONY: cover
cover:
	go tool cover \
		-html \
		cover/cover.out

.PHONY: vet
vet:
	go vet -v $(PKGS)

.PHONY: lint
lint:
	golangci-lint run -v ./...
