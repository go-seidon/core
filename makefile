
default: help

.PHONY: help
help:
	@echo 'goseidon'
	@echo 'usage: make [target] ...'

.PHONY: install-tool
install-tool:
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen

.PHONY: install-dependency
install-dependency:
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: clean-dependency
clean-dependency:
	rm -f go.sum
	rm -rf vendor
	go clean -modcache

.PHONY: install
install:
	go install -v ./...

.PHONY: test
test:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out | grep ^total:

.PHONY: generate-mock
generate-mock:
	mockgen -source goseidon.go -destination=goseidon_mock.go -package=goseidon
	mockgen -source pkg/aws-s3/storage.go -destination=pkg/aws-s3/storage_mock.go -package=aws_s3
	mockgen -source internal/io/client.go -destination=internal/io/client_mock.go -package=io

.PHONY: run-example
run-example:
	go run example/main.go
