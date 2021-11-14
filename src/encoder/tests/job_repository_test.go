

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryDb_InsertNewJobReturnSucess(t *testing.T) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), uuid.NewV4().String(), "test_path", time.Now())

	db := database.NewDbTest()
	defer db.Close()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, job.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)

}

func TestJobRepositoryDb_UpdateJobReturnJob(t *testing.T) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), uuid.NewV4().String(), "test_path", time.Now())
	db := database.NewDbTest()
	defer db.Close()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "Complete"
	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, job.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
	require.Equal(t, j.Status, job.Status)

}
