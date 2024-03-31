run:
	go run .

go-licenses:
	bash hack/go-licenses.sh

gocyclo:
	bash hack/gocyclo.sh

misspell:
	bash hack/misspell.sh

test:
	go test ./... --failfast -v

cover:
	go test ./... --failfast -coverprofile /tmp/coverage.out
	go tool cover -func /tmp/coverage.out | tail -1

checks:
	bash hack/checks.sh
