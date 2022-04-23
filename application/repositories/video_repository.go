package repositories

import (
	uuid "github.com/satori/go.uuid"
	"github.com/video-encoder/domain"
	"gorm.io/gorm"
)

type VideoRepositoryDB struct {
	DB *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{db}
}

func (ref VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	if err := ref.DB.Create(video).Error; err != nil {
		return nil, err
	}

	return video, nil
}
func (ref VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video

	if err := ref.DB.Where(&domain.Video{ID: id}).First(&video).Error; err != nil {
		return nil, err
	}

	return &video, nil
}
