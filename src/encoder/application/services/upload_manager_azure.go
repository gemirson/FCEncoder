package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type VideoUploadAzure struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Erros        []string
}

func NewVideoUploadAzure() *VideoUploadAzure {
	return &VideoUploadAzure{}
}

func (vu *VideoUploadAzure) UploadObject(objectpath string, client *azblob.ContainerURL, ctx context.Context) error {

	path := strings.Split(objectpath, os.Getenv("localStoragePath")+"/")

	f, err := ReadFile(objectpath)
	//	f, err := os.Open(objectpath)

	if err != nil {
		return err
	}

	//defer f.Close()

	blockBlobURL := client.NewBlockBlobURL(path[1]) // Blob names can be mixed case

	o := azblob.UploadToBlockBlobOptions{BlobHTTPHeaders: azblob.BlobHTTPHeaders{ContentType: "text/plain"}}

	_, err = azblob.UploadBufferToBlockBlob(ctx, f, blockBlobURL, o)

	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func ReadFile(filePath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

func (vu *VideoUploadAzure) loadPaths() error {

	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUploadAzure) ProcessUpload(concurrency int, doneUpload chan string, uploadClient *azblob.ContainerURL) error {

	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()

	if err != nil {
		return err
	}
	ctx := vu.getClientContextUpload()

	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)

	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x

		}
		close(in)

	}()

	for r := range returnChannel {

		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil

}

func (vu *VideoUploadAzure) uploadWorker(in chan int, returnChan chan string, uploadClient *azblob.ContainerURL, ctx context.Context) {

	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Erros = append(vu.Erros, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v ", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func (vu *VideoUploadAzure) getClientContextUpload() context.Context {

	ctx := context.Background()

	return ctx
}
