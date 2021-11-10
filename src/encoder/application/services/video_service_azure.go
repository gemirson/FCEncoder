package services

import (
	"bytes"
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type VideoServiceAzure struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoServiceAzure() VideoServiceAzure {
	return VideoServiceAzure{}
}

func (v *VideoServiceAzure) Download(bucketName string) error {

	ctx := context.Background()

	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
		return err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, bucketName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create the container
	fmt.Printf("Creating a container named %s\n", bucketName)

	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)

	if err != nil {
		return err
	}

	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := containerURL.NewBlockBlobURL(v.Video.FilePath) // Blob names can be mixed case

	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	downloadedData := &bytes.Buffer{}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	downloadedData.ReadFrom(bodyStream)

	if err != nil {
		log.Fatal(err)
	}
	defer bodyStream.Close()

	body, err := ioutil.ReadAll(bodyStream)

	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}
	_, err = f.Write(body)

	if err != nil {
		return err
	}

	log.Printf("video %v has been stored", v.Video.ID)

	return nil

}
