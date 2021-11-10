package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	service    azblob.ServiceClient
	serviceURL = os.Getenv("AZURE_STORAGE_SERVICE_URL")
	ctx        = context.Background()
)

func init() {

	fmt.Printf("Storage service URL: %s\n", serviceURL)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Print("Cannot create credentials")
		panic("Cannot create credentials")
	}

	service, err = azblob.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		fmt.Print("Cannot connect to service")
		panic("Cannot connect to service")
	}
}

func task1(containerName, blobName, blobData string) {
	// Write code that uploads the string blobData to a blob named blobName in the container containerName

	container := service.NewContainerClient(containerName)
	_, err := container.Create(ctx, nil)

	// Should the Create method return a different error type? How do I determine that the error is because the container
	// already exists? All I have is the error string since the error instance only implements the error interface
	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("ErrorCode=%s", azblob.StorageErrorCodeContainerAlreadyExists)) {
			fmt.Printf("Error: %s", err.Error())
			panic("Cannot access the container")
		}
	}

	blockBlob := container.NewBlockBlobClient(blobName)

	// Unlike containers, creating a block blob with the name of one that already exists is not an error
	// If users want to add to a blob, they should use an append blob client

	// Note that I can pass nil as the options parameter

	_, err = blockBlob.Upload(ctx, streaming.NopCloser(strings.NewReader(blobData)), nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

}

func task2(fileName, containerName, blobName, contentType string) {

	// Write code that uploads the file filename to blob blobname in container containername
	// The content in the file is of type contentType
	container := service.NewContainerClient(containerName)
	largeFile := container.NewBlockBlobClient(blobName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer file.Close()

	_, err = container.Create(ctx, nil)
	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("ErrorCode=%s", azblob.StorageErrorCodeContainerAlreadyExists)) {
			fmt.Printf("Error: %s", err.Error())
		}
	}

	conType := new(azblob.BlobHTTPHeaders)
	conType.BlobContentType = &contentType

	// Note the singular Option in the name
	options := azblob.HighLevelUploadToBlockBlobOption{HTTPHeaders: conType}
	// Why is the method named UploadFileToBlockBlob when the client is a BlockBlobClient
	_, err = largeFile.UploadFileToBlockBlob(ctx, file, options)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

}

func task3(fileName, containerName, blobName, contentType string) {
	// Modify the code that you wrote in task 2 so that the upload operation is cancelled
	// if it takes longer than 5 seconds

	container := service.NewContainerClient(containerName)
	largeFile := container.NewBlockBlobClient(blobName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer file.Close()

	_, err = container.Create(ctx, nil)
	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("ErrorCode=%s", azblob.StorageErrorCodeContainerAlreadyExists)) {
			fmt.Printf("Error: %s", err.Error())
		}
	}

	conType := new(azblob.BlobHTTPHeaders)
	conType.BlobContentType = &contentType

	// Note the singular Option in the name
	options := azblob.HighLevelUploadToBlockBlobOption{HTTPHeaders: conType}

	// Setting up the upload to time out with the context object
	// Do they expect to do this in the context object or in the options?
	//
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// How to cancel this if upload takes a long time - do go developers expect to use the context object?

	_, err = largeFile.UploadFileToBlockBlob(ctx, file, options)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

func task4(containerName, blobName string) {
	// Write code that retrieves the content of the blob named blobName in the container containerName

	container := service.NewContainerClient(containerName)
	blob := container.NewBlobClient(blobName)
	response, err := blob.Download(ctx, nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("Cannot access the blob")
	}

	buf := new(bytes.Buffer)
	// Cannot pass nil to the options - need to create an empty RetryReaderOptions
	// Is it clear when the data is actually downloaded? It's interesting that response.Body
	// takes a retryreaderoptions which lets us set maxretries - does that mean that download
	// hasn't actually downloaded the data? In other
	n, err := buf.ReadFrom(response.Body(azblob.RetryReaderOptions{}))
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	newStr := buf.String()[0:80]
	fmt.Printf("Downloaded %s bytes: %s\n", fmt.Sprint(n), newStr)

}

func task5(containerName, blobName string, numRetries int) {
	// Write code that retrieves the content of the blob named blobName in the container containerName
	// and that is configured to retry retrieving the content up to numRetries times should any attempt fail

	container := service.NewContainerClient(containerName)
	blob := container.NewBlobClient(blobName)
	response, err := blob.Download(ctx, nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	buf := new(bytes.Buffer)

	// Testing where they think the download actually happens. Do they look at the options on the
	// download function or on the Response.Body
	// What are their thoughts about why they can pass nil as the parameter for the download options
	// but not the body options?
	n, err := buf.ReadFrom(response.Body(azblob.RetryReaderOptions{MaxRetryRequests: numRetries}))
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	newStr := buf.String()[0:80]
	fmt.Printf("Downloaded %s bytes: %s\n", fmt.Sprint(n), newStr)

}

func task6() {
	// Write code that prints out the name of every container in the storage account
	// Optional additional task - only print out the names of any container if the name begins with "container"
	// options can be nil
	options := new(azblob.ListContainersOptions)
	prefix := "container"
	options.Prefix = &prefix
	pager := service.ListContainers(options)

	// Is it clear that NextPage needs to be called first?
	for pager.NextPage(ctx) {
		for index, element := range pager.PageResponse().ContainerItems {
			fmt.Println("At index", index, "value is", *element.Name)
		}
	}

}

func task7(containerName, logName string) {
	// Write code that writes a series of log entries to a log

	container := service.NewContainerClient(containerName)
	log := container.NewAppendBlobClient(logName)
	for i := 0; i < 100; i++ {
		event := "Event happened at: " + time.Now().String()
		log.AppendBlock(ctx, streaming.NopCloser(strings.NewReader(event)), nil)

	}
}

func task8() {
	// Write code that deletes all of the containers
	pager := service.ListContainers(nil)

	// Is it clear that NextPage needs to be called first?
	for pager.NextPage(ctx) {
		for _, element := range pager.PageResponse().ContainerItems {
			service.DeleteContainer(ctx, *element.Name, nil)
		}
	}
}

func main() {

	task1("helloworldcontainer", "helloworldblob", "Hello World at "+time.Now().String())
	// task2("enwik9.pmd", "testcontainer", "largefile", "text/xml")
	// task3("enwik9.pmd", "testcontainer", "largefile", "text/xml")
	// task4("task4", "blob.txt")
	// task5("testcontainer", "largefile", 5)
	// task6()
	// task7("logcontainer", "log")
	// task8()

	// Use this code to create the containers before every session
	// // Create 1000 containers
	// // for i := 0; i < 1000; i++ {
	// // 	_, err := service.CreateContainer(ctx, "container"+strconv.Itoa(i), nil)
	// // 	if err != nil {
	// // 		fmt.Printf("Error: %s", err.Error())
	// // 	}
	// // }
	//
	// or
	// seq 1 100 | while read d; do az storage container create -n container${d}; done

	fmt.Println("End of program")
}
