VERSION := v0.2.0

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
