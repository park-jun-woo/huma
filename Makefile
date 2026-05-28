VERSION := v0.1.2

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
