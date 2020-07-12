PKGS := $(shell go list ./... | grep -v /vendor)

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
