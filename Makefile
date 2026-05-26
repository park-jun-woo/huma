VERSION := v0.1.0

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
