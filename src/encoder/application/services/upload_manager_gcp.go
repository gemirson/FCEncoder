package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUploadCGP struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Erros        []string
}

func NewVideoUpload() *VideoUploadCGP {
	return &VideoUploadCGP{}
}

func (vu *VideoUploadCGP) UploadObject(objectpath string, client *storage.Client, ctx context.Context) error {

	path := strings.Split(objectpath, os.Getenv("localStoragePath")+"/")
	f, err := os.Open(objectpath)

	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

func (vu *VideoUploadCGP) loadPaths() error {

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

func (vu *VideoUploadCGP) ProcessUpload(concurrency int, doneUpload chan string) error {

	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()

	if err != nil {
		return err
	}
	uploadClient, ctx, err := getClientUpload()

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

func (vu *VideoUploadCGP) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {

	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Erros = append(vu.Erros, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v ", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload complete"
	returnChan <- "upload complete"
}

func getClientUpload() (*storage.Client, context.Context, error) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
