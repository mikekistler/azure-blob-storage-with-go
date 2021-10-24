package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

const (
	localFileName = "file.txt"
)

var (
	ctx             context.Context
	serviceClient   azblob.ServiceClient
	containerClient azblob.ContainerClient
)

func handle(err error) {

	if err != nil {
		fmt.Println("Error! " + err.Error())
	}

}

func createServiceClient() {

	// Obtain the storage account connection string from the environment
	connStr := os.Getenv("STORAGE_CONNECTION_STRING")

	// Create the service client
	var err error
	serviceClient, err = azblob.NewServiceClientFromConnectionString(connStr, nil)
	handle(err)

}

func createContainerClient() {

	// Create a unique name for the container
	containerName := uuid.New().String()

	// Create the container client
	containerClient = serviceClient.NewContainerClient(containerName)

}

func createContainer() {

	// Create the container
	_, err := containerClient.Create(ctx, nil)
	handle(err)

}

func uploadFile() {

	// Read file contents into byte variable
	data, err := ioutil.ReadFile(localFileName)
	handle(err)

	blockBlob := containerClient.NewBlockBlobClient(localFileName)
	_, err = blockBlob.Upload(ctx, streaming.NopCloser(bytes.NewReader(data)), nil)
	handle(err)

}

func listBlobs() {

	pager := containerClient.ListBlobsFlat(nil)

	fmt.Println("Blobs in container")
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			fmt.Println("\t" + *v.Name)
		}
	}

}

func main() {
	fmt.Println("Azure Blob storage quick start sample\n")

	// All operations in the Azure Storage Blob SDK for Go operate on a context.Context, allowing you to control cancellation/timeout.
	ctx = context.Background() // This example has no expiry.

	// Create service client
	createServiceClient()

	// Create container client
	createContainerClient()

	// Create container
	createContainer()

	// Upload file
	uploadFile()

	// List blobs in the container
	listBlobs()
}
