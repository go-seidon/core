package main

import (
	"fmt"
	"os"

	goseidon "github.com/go-seidon/core"
	aws_s3 "github.com/go-seidon/core/pkg/aws-s3"
)

func main() {
	fmt.Println("+==========+")
	fmt.Println("| Goseidon |")
	fmt.Println("+==========+")
	s3()
}

func s3() {
	fmt.Println()
	fmt.Println("Trying Goseidon Storage with S3 provider")
	fmt.Println("==================================")

	region := os.Getenv("AWS_REGION")
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	cr, err := aws_s3.NewAwsS3Credential(region, accessKeyId, secretAccessKey)
	if err != nil {
		panic(err)
	}

	cl, err := aws_s3.NewAwsS3Client(cr)
	if err != nil {
		panic(err)
	}

	op, err := aws_s3.NewAwsS3Option(bucketName)
	if err != nil {
		panic(err)
	}

	storage, err := aws_s3.NewAwsS3Storage(cl, op)
	if err != nil {
		panic(err)
	}

	osFile, err := os.Open("example/dolphin.jpg")
	if err != nil {
		panic(err)
	}
	defer osFile.Close()

	fileInfo, _ := osFile.Stat()
	var fileSize int64 = fileInfo.Size()
	fileData := make([]byte, fileSize)
	osFile.Read(fileData)

	uploadRes, err := storage.UploadFile(goseidon.UploadFileParam{
		FileName: fileInfo.Name(),
		FileData: fileData,
		FileSize: fileSize,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Upload Result => ", uploadRes)

	retrieveRes, err := storage.RetrieveFile(goseidon.RetrieveFileParam{
		Id: fileInfo.Name(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieve Result => ", retrieveRes)

	fmt.Println("==================================")
	fmt.Println("Finish trying Goseidon Storage with S3 provider")
	fmt.Println()

	fmt.Println("Don't forget to delete the uploaded file in your bucket!")
	fmt.Println("Press any key to continue...")
	fmt.Scanln()
}
