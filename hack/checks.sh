#!/bin/bash
cd $(dirname $0)/../

set -xeuo pipefail
go mod tidy
go fmt ./...
go vet ./...

bash hack/go-licenses.sh
bash hack/go-test.sh
bash hack/gocyclo.sh
bash hack/golangci-lint.sh
bash hack/misspell.sh
bash hack/staticcheck.sh

which goimports || go install golang.org/x/tools/cmd/goimports@latest
goimports -local -v -w .

echo "✔️ OK"
