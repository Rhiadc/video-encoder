package repositories

import (
	uuid "github.com/satori/go.uuid"
	"github.com/video-encoder/domain"
	"gorm.io/gorm"
)

type JobRepositoryDB struct {
	DB *gorm.DB
}

func NewJobRepositoryDB(db *gorm.DB) *JobRepositoryDB {
	return &JobRepositoryDB{db}
}

func (ref JobRepositoryDB) Insert(job *domain.Job) (*domain.Job, error) {
	if job.ID == "" {
		job.ID = uuid.NewV4().String()
	}

	if err := ref.DB.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (ref JobRepositoryDB) Find(id string) (*domain.Job, error) {
	var job domain.Job

	if err := ref.DB.Joins("Video").First(&job, id).Error; err != nil {
		return nil, err
	}

	return &job, nil
}
