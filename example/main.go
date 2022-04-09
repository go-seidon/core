package main

import (
	"bufio"
	"fmt"
	"os"

	goseidon "github.com/go-seidon/core"
	aws_s3 "github.com/go-seidon/core/pkg/aws-s3"
	"github.com/go-seidon/core/pkg/local"
)

func main() {
	fmt.Println("+==========+")
	fmt.Println("| Goseidon |")
	fmt.Println("+==========+\n\r")

	fmt.Println("[1] Local Storage")
	fmt.Println("[2] AWS S3")
	fmt.Print("Choose your storage provider: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	menu := scanner.Text()

	fmt.Printf("Choosen menu: %s \n\r", menu)

	switch menu {
	case "1":
		TryLocal()
	case "2":
		TryS3()
	default:
		fmt.Println()
		fmt.Println("Invalid storage provider")
		fmt.Println("Please choose between [1, 2]")
	}

}

func TryLocal() {
	fmt.Println()
	fmt.Println("Trying Goseidon Storage with Local provider")
	fmt.Println("==================================")

	config, err := local.NewLocalConfig("storage")
	if err != nil {
		panic(err)
	}

	storage, err := local.NewLocalStorage(config)
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
	fmt.Println("Finish trying Goseidon Storage with Local provider")
	fmt.Println()

	fmt.Printf("Don't forget to delete the uploaded file in: %s \n\r", config.StorageDir)
	fmt.Println("Press any key to continue...")

	fmt.Scanln()
}

func TryS3() {
	fmt.Println()
	fmt.Println("Trying Goseidon Storage with S3 provider")
	fmt.Println("==================================")

	region := os.Getenv("AWS_REGION")
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	cfg := &aws_s3.AwsS3Config{
		Credential: &aws_s3.AwsS3Credential{
			Region:          region,
			AccessKeyId:     accessKeyId,
			SecretAccessKey: secretAccessKey,
		},
		BucketName: bucketName,
	}
	storage, err := aws_s3.NewAwsS3Storage(cfg)
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

	fmt.Printf("Don't forget to delete the uploaded file in: %s\n\r", cfg.BucketName)
	fmt.Println("Press any key to continue...")
	fmt.Scanln()
}
