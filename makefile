
default: help

.PHONY: help
help:
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@echo 'generate-mock'

install-tool:
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen

install-dependency:
	go mod tidy
	go mod verify
	go mod vendor

clean-dependency:
	rm -f go.sum
	rm -rf vendor
	go clean -modcache

install:
	go install -v ./...

run:
	go run main.go

test:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out | grep ^total:

generate-mock:
	mockgen -source goseidon.go -destination=goseidon_mock.go -package=goseidon
	mockgen -source pkg/aws-s3/storage.go -destination=pkg/aws-s3/storage_mock.go -package=aws_s3
