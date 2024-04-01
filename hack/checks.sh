#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go mod tidy
go fmt ./...
go vet ./...

bash hack/go-licenses.sh
bash hack/gocyclo.sh
bash hack/misspell.sh
bash hack/staticcheck.sh

which goimports || go install golang.org/x/tools/cmd/goimports@latest
goimports -local -v -w .

which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --timeout 5m

which go-licenses || go install github.com/google/go-licenses@v1.6.0
go-licenses check ./...

echo "✔️ OK"