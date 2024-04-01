run:
	go run .

go-licenses:
	bash hack/go-licenses.sh
gocyclo:
	bash hack/gocyclo.sh
golangci-lint:
	bash hack/golangci-lint.sh
misspell:
	bash hack/misspell.sh
staticcheck:
	bash hack/staticcheck.sh

test:
	go test ./... --failfast -v

cover:
	go test ./... --failfast -coverprofile /tmp/coverage.out
	go tool cover -func /tmp/coverage.out | tail -1

checks:
	bash hack/checks.sh

build:
	CGO_ENABLED=0 GOOS=linux go build -C pkg -o /tmp/runbox

docker:
	docker build -t runbox .
