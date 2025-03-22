.PHONY: build

# Show this help.
# Usage: `make help`
help:
	@cat $(MAKEFILE_LIST) | docker run --rm -i xanders/make-help

# Build the project.
# By default, it builds for the current platform.
# Usage: `make build`
# You can specify a different platform by setting the OS and ARCH variables.
# Usage: `make build OS=windows ARCH=x86_64`
build:
	scripts/build.sh $(OS) $(ARCH)

# Run project from source code for the current OS and architecture.
# Usage: `make run`
run:
	go run cmd/algorewards/main.go

# Lint the project source code.
# Usage: `make lint`
lint:
	golangci-lint run

# Release the project.
# Usage: `make release`
release:
	$(MAKE) build OS=all ARCH=all
	goreleaser release --clean

# Generate a snapshot release.
# Usage: `make snapshot`
snapshot:
	goreleaser release --clean --snapshot

# Install release tools for macOS.
# Usage: `make install-releasetools`
install-releasetools:
	brew install goreleaser/tap/goreleaser-pro
	brew install msitools