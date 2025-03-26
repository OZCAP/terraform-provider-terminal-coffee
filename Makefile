default: build

HOSTNAME=registry.terraform.io
NAMESPACE=OZCAP
NAME=terminal-coffee
VERSION=0.1.0
OS_ARCH=darwin_amd64

.PHONY: build
build:
	go build -o terraform-provider-${NAME} ./main

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp terraform-provider-${NAME} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f terraform-provider-${NAME}