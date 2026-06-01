VERSION := v0.2.1

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" .
