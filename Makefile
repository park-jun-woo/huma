VERSION := v0.1.4

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
