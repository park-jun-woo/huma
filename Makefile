VERSION := v0.1.3

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
