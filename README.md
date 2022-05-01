# goseidon

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_core&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-seidon_core)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_core&metric=coverage)](https://sonarcloud.io/summary/new_code?id=go-seidon_core)

File uploader library for uploading to multiple storage provider

Current support:
- `local file upload`
- `aws s3`
- `g-cloud storage`

Upcoming support:
- `alicloud oss`

## Doc
See [code example](example/main.go) for the moment 😁

## Todo
1. Add `RetrievedAt` in `RetrieveFileResult`
2. Add `FileId` in `UploadFileParam` & `UploadFileResult` (replace FileName as file identifier)
3. Refactor storage factory to receive `client` param
4. Remove `AwsCredential` (flat is better than fat)
