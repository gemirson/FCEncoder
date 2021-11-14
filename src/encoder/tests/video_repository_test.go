package tests_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDb_InsertNewVideoReturnSucess(t *testing.T) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), uuid.NewV4().String(), "test_path", time.Now())

	db := database.NewDbTest()
	defer db.Close()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, video.ID)
	require.NotNil(t, v)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)

}
