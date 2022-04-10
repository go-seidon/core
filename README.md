# goseidon

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_core&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-seidon_core)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_core&metric=coverage)](https://sonarcloud.io/summary/new_code?id=go-seidon_core)

File uploader library for uploading to multiple storage provider

Current support:
- `aws s3`
- `local file upload`

Upcoming support:
- `google storage`
- `alicloud oss`

## Doc
See [code example](example/main.go) for the moment 😁

## Todo
1. aws s3 storage: simplify `NewAwsS3Storage` parameter (do we need `AWSS3Option`?)
2. local storage: simplify `NewLocalStorage` parameter (seperate `io.FileManagerService` for test and client usage)
3. local storage: should create `storageDir` if not exists
