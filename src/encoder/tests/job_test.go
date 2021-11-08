package tests_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobEntity_CreateNewJob_returnJob(t *testing.T) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), uuid.NewV4().String(), "test_path", time.Now())
	job, err := domain.NewJob("test_path", "test_Converted", video)

	require.NotNil(t, job)
	require.Nil(t, err)

}
