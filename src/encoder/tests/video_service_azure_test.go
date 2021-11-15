package tests_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceAzure_DownloadVideoReturnSucess(t *testing.T) {

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

	err = videoService.Finish()
	require.Nil(t, err)

}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), uuid.NewV4().String(), "scorpions-wind-of-change.mp4", time.Now())

	db := database.NewDbTest()
	defer db.Close()

	repo := repositories.VideoRepositoryDb{Db: db}
	return video, repo
}
