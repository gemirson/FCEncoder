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
	"os/exec"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type VideoServiceAzure struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
	AzureContainer  *azblob.ContainerURL
}

func NewVideoServiceAzure() VideoServiceAzure {
	return VideoServiceAzure{}
}

func (v *VideoServiceAzure) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, v.ExtractedPathTarget())
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, v.ExtractedPathDirectoryTarget())
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoServiceAzure) Download(bucketName string) error {

	ctx := context.Background()

	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := AccountInfo()

	ValidateAccountInformation(accountName, accountKey)

	// Create a default request pipeline using your storage account name and account key.
	// From the Azure portal, get your storage account blob service URL endpoint.
	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL, err := v.GenerateAndConfiguredContainer(accountName, accountKey, bucketName)

	if err != nil {
		return err
	}

	// Create the container
	fmt.Printf("Reading a container named %s\n", bucketName)
	v.AzureContainer = &containerURL
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := containerURL.NewBlockBlobURL(v.Video.FilePath) // Blob names can be mixed case

	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	downloadedData := &bytes.Buffer{}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	downloadedData.ReadFrom(bodyStream)
	defer bodyStream.Close()

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(downloadedData)

	if err != nil {
		return err
	}

	f, err := v.ExtractedLocalPath()

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

func ValidateAccountInformation(accountName string, accountKey string) {
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}
}

func (v *VideoServiceAzure) GenerateAndConfiguredContainer(accountName string, accountKey string, bucketName string) (azblob.ContainerURL, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
		return azblob.ContainerURL{}, err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, bucketName))

	containerURL := azblob.NewContainerURL(*URL, p)

	return containerURL, nil
}

func (v *VideoServiceAzure) ExtractedLocalPath() (*os.File, error) {

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return nil, err
	}
	return f, nil
}

func (v *VideoServiceAzure) Fragment() error {

	err := os.Mkdir(v.ExtractedPathDirectoryTarget(), os.ModePerm)
	if err != nil {
		return err
	}

	source := v.ExtractedPathSource()
	target := v.ExtractedPathTarget()

	cmd := exec.Command("mp4fragment", source, target)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoServiceAzure) Finish() error {

	err := os.Remove(v.ExtractedPathSource())
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+".mp4")
		return err
	}

	err = os.Remove(v.ExtractedPathTarget())
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+".frag")
		return err
	}

	err = os.RemoveAll(v.ExtractedPathDirectoryTarget())
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+".frag")
		return err
	}

	log.Println("files have been removed !", v.Video.ID)

	return nil
}

func AccountInfo() (string, string) {

	return os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
}

func (v *VideoServiceAzure) ExtractedPathSource() string {
	return os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
}

func (v *VideoServiceAzure) ExtractedPathTarget() string {
	return os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"
}

func (v *VideoServiceAzure) ExtractedPathDirectoryTarget() string {
	return os.Getenv("localStoragePath") + "/" + v.Video.ID
}

func (v *VideoServiceAzure) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)

	if err != nil {
		return err
	}

	return nil
}
func printOutput(out []byte) {

	if len(out) > 0 {
		log.Printf("=====================> Output: %s\n", string(out))
	}

}
