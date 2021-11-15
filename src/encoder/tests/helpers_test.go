package tests_test

import (
	"encoder/framework/helpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
  				"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
  				"file_path": "convite.mp4",
  				"status": "pending"
			}`

	err := helpers.IsJson(json)
	require.Nil(t, err)

	json = `wes`
	err = helpers.IsJson(json)
	require.Error(t, err)
}
