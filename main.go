package main

import (
	"context"
	"fmt"
	"time"
)

var (
	ctx = context.Background()
)

func task1(containerName, blobName, blobData string) {

	// Write code that uploads the string blobData to a blob named blobName in the container containerName

}

func task2(fileName, containerName, blobName, contentType string) {

	// Write code that uploads the file filename to blob blobname in container containername
	// The content in the file is of type contentType

}

func task3(fileName, containerName, blobName, contentType string) {

	// Modify the code that you wrote in task 2 so that the upload operation is cancelled
	// if it takes longer than 5 seconds

}

func task4(containerName, blobName string) {

	// Write code that retrieves the content of the blob named blobName in the container containerName

}

func task5(containerName, blobName string, numRetries int) {

	// Write code that retrieves the content of the blob named blobName in the container containerName
	// and that is configured to retry retrieving the content up to numRetries times should any attempt fail

}

func task6() {

	// Write code that prints out the name of every container in the storage account
	// Optional additional task - only print out the names of any container if the name begins with "container"

}

func task7(containerName, logName string) {

	// Write code that writes a series of log entries to a log

}

func task8() {

	// Write code that deletes all of the containers

}

func main() {

	task1("helloworldcontainer", "helloworldblob", "Hello World at "+time.Now().String())
	task2("enwik9.pmd", "testcontainer", "largefile", "text/xml")
	task3("enwik9.pmd", "testcontainer", "largefile", "text/xml")
	task4("task4", "blob.txt")
	task5("testcontainer", "largefile", 5)
	task6()
	task7("logcontainer", "log")
	task8()

	fmt.Println("End of program")
}
