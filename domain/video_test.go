package domain_test

import (
	"github.com/stretchr/testify/require"
	"github.com/video-encoder/domain"
	"testing"
	"time"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.Video{}
	require.Error(t, video.Validate())

}

func TestVideoIDIsNotAUUID(t *testing.T) {
	video := domain.Video{}
	video.ID = "not_a_uuid"
	video.FilePath = "some_path"
	video.ResourceID = "some_resource_id"
	video.CreatedAt = time.Now()
	require.Error(t, video.Validate())
}
