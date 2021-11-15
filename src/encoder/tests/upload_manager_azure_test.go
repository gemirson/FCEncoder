package tests_test

import (
	"encoder/application/services"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestUploadManagerAzure_UploadVideoReturnSucess(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoServiceAzure()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("codeflix")

	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUploadAzure()
	videoUpload.OutputBucket = "codeflix"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(50, doneUpload, videoService.AzureContainer)

	result := <-doneUpload

	require.Equal(t, result, "upload completed")

}
