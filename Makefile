default: build

HOSTNAME=registry.terraform.io
NAMESPACE=OZCAP
NAME=terminal-coffee
VERSION=1.0.8
OS_ARCH=darwin_amd64
BINARY=terraform-provider-${NAME}
SIGNING_KEY=$(shell git config --get user.signingkey)

PLATFORMS=darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64
RELEASE_DIR=releases

.PHONY: build
build:
	go build -o ${BINARY} ./main

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f ${BINARY}
	rm -rf ${RELEASE_DIR}

.PHONY: release
release: clean
	mkdir -p ${RELEASE_DIR}
	@echo "Building for platforms: ${PLATFORMS}"
	@for platform in ${PLATFORMS}; do \
		os=$$(echo $$platform | cut -d_ -f1); \
		arch=$$(echo $$platform | cut -d_ -f2); \
		extension=""; \
		if [ "$$os" = "windows" ]; then extension=".exe"; fi; \
		echo "Building for $${os}_$${arch}"; \
		GOOS=$$os GOARCH=$$arch go build -o ${RELEASE_DIR}/${BINARY}_v${VERSION}_$${os}_$${arch}$$extension ./main; \
		(cd ${RELEASE_DIR} && zip ${BINARY}_v${VERSION}_$${os}_$${arch}.zip ${BINARY}_v${VERSION}_$${os}_$${arch}$$extension); \
	done
	(cd ${RELEASE_DIR} && shasum -a 256 *.zip > ${BINARY}_v${VERSION}_SHA256SUMS)
	(cd ${RELEASE_DIR} && export GPG_TTY=$$(tty) && gpg --detach-sign ${BINARY}_v${VERSION}_SHA256SUMS || echo "WARNING: GPG signing failed, but continuing build")
	@echo "Release files created in ${RELEASE_DIR}"
	@echo "To create a GitHub release and tag, run:"
	@echo "  git tag -s v${VERSION} -u ${SIGNING_KEY} -m \"Release v${VERSION}\""
	@echo "  git push origin v${VERSION}"
	@echo "  gh release create v${VERSION} --title \"v${VERSION}\" --notes \"Release notes\" ${RELEASE_DIR}/*.zip ${RELEASE_DIR}/${BINARY}_v${VERSION}_SHA256SUMS ${RELEASE_DIR}/${BINARY}_v${VERSION}_SHA256SUMS.sig"

.PHONY: release-tag
release-tag:
	git tag -s v${VERSION} -u ${SIGNING_KEY} -m "Release v${VERSION}"
	git push origin v${VERSION}

.PHONY: github-release
github-release:
	gh release create v${VERSION} --title "v${VERSION}" --notes "Release v${VERSION}" ${RELEASE_DIR}/*.zip ${RELEASE_DIR}/${BINARY}_v${VERSION}_SHA256SUMS ${RELEASE_DIR}/${BINARY}_v${VERSION}_SHA256SUMS.sig