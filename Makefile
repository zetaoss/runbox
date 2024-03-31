run:
	go run .

test:
	go test ./... --failfast -v

cover:
	go test ./... --failfast -coverprofile /tmp/coverage.out
	go tool cover -func /tmp/coverage.out | tail -1

checks: test
