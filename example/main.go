package main

import (
	"bufio"
	"context"
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

	var storage goseidon.Storage

	switch menu {
	case "1":
		storage = MustCreateLocalStorage()
	case "2":
		storage = MustCreateAwsS3Storage()
	default:
		fmt.Println()
		fmt.Println("Invalid storage provider")
		fmt.Println("Please choose between [1, 2]")
		return
	}

	fmt.Println()
	fmt.Println("Trying Goseidon Storage")
	fmt.Println("==================================")

	osFile, err := os.Open("example/dolphin.jpg")
	if err != nil {
		panic(err)
	}
	defer osFile.Close()

	fileInfo, _ := osFile.Stat()
	var fileSize int64 = fileInfo.Size()
	fileData := make([]byte, fileSize)
	osFile.Read(fileData)

	ctx := context.Background()
	uploadRes, err := storage.UploadFile(ctx, goseidon.UploadFileParam{
		FileName: fileInfo.Name(),
		FileData: fileData,
		FileSize: fileSize,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Upload Result => ", uploadRes)

	retrieveRes, err := storage.RetrieveFile(ctx, goseidon.RetrieveFileParam{
		Id: fileInfo.Name(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieve Result => ", retrieveRes)

	fmt.Println("==================================")
	fmt.Println("Please check the uploaded file")
	fmt.Println("==================================")
	fmt.Scanln()

	fmt.Println("Press any key to delete the uploaded file...")
	fmt.Scanln()

	deleteRes, err := storage.DeleteFile(ctx, goseidon.DeleteFileParam{
		Id: fileInfo.Name(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("File deleted at: %s \n\r", deleteRes.DeletedAt.Local())
	fmt.Println("==================================")
	fmt.Println()

	fmt.Println("Press any key to continue...")

	fmt.Scanln()
}

func MustCreateLocalStorage() goseidon.Storage {
	config, err := local.NewLocalConfig("storage")
	if err != nil {
		panic(err)
	}

	storage, err := local.NewLocalStorage(config)
	if err != nil {
		panic(err)
	}

	return storage
}

func MustCreateAwsS3Storage() goseidon.Storage {
	region := os.Getenv("AWS_REGION")
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	cfg, _ := aws_s3.NewAwsS3Config(
		region,
		accessKeyId,
		secretAccessKey,
		bucketName,
	)

	storage, err := aws_s3.NewAwsS3Storage(cfg)
	if err != nil {
		panic(err)
	}

	return storage
}
