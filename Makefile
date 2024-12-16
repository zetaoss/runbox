GOLANGCI_LINT_VER := v1.59.1
GO_LICENSES_VER := v1.6.0
IMG := jmnote/runbox:runbox

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./... --failfast -v

.PHONY: cover
cover:
	go test ./... --failfast -coverprofile /tmp/coverage.out
	go tool cover -func /tmp/coverage.out | tail -1

.PHONY: lint
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER) || true
	$(shell go env GOPATH)/bin/golangci-lint run

.PHONY: checks
checks: test lint

.PHONY: docker-build
docker-build:
	docker build -t $(IMG) .

.PHONY: docker-push
docker-push:
	docker push $(IMG)

.PHONY: docker
docker: docker-build docker-push
