package tests_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoEntity_ValidateIfVideoIsEmpty(t *testing.T) {

	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)

}

func TestVideoEntity_ValidateVideoIdIsNotAUuid(t *testing.T) {

	video := domain.CreateInstanceVideo("teste", "teste", "test_path", time.Now())
	err := video.Validate()

	require.Error(t, err)

}

func TestVideoEntity_ValidateVideoIdIsAUuid(t *testing.T) {

	video := domain.CreateInstanceVideo(uuid.NewV4().String(), "teste", "test_path", time.Now())
	err := video.Validate()
	require.Nil(t, err)

}
