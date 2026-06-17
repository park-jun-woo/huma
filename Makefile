VERSION := v0.3.0

.PHONY: install build test vet
install:
	go install -ldflags "-X main.Version=$(VERSION)" .

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o huma .

vet:
	go vet ./internal/...

test:
	go test ./internal/...
