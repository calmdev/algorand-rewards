FROM golang:1.24

RUN apt-get update && \
	# Go tools:
	# - https://github.com/golang/tools
	# - https://github.com/golang/vscode-go/blob/HEAD/docs/tools.md
	go install -v golang.org/x/tools/gopls@latest && \
	go install github.com/cweill/gotests/gotests@latest && \
	go install github.com/fatih/gomodifytags@latest && \
	go install github.com/josharian/impl@latest && \
	go install github.com/haya14busa/goplay/cmd/goplay@latest && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/cmd/godoc@latest && \
	go install gotest.tools/gotestsum@latest && \
	# Go vulnerability checker:
	# https://go.dev/security/vuln/
	go install golang.org/x/vuln/cmd/govulncheck@latest && \
	# Fyne for GUI development:
	# - https://github.com/fyne-io/fyne
	go install fyne.io/fyne/v2/cmd/fyne@latest && \
	go install github.com/fyne-io/fyne-cross@latest && \
	# golangci-lint for linting:
	# - https://golangci-lint.run
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 && \
	# Graphviz for visualizing go test pprof profiles:
	# - https://graphviz.org
	apt-get -y install graphviz && \
	# Docker for containerization:
	# - https://www.docker.com
	curl -fsSL https://get.docker.com -o get-docker.sh && \
	sh get-docker.sh

# ADD . /workspace
# WORKDIR /workspace

CMD ["sleep", "infinity"]