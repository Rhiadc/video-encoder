package domain

import (
	"github.com/asaskevich/govalidator"
	"time"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceID string    `valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"-" valid:"notnull"`
	CreatedAt  time.Time `json:"-" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (ref *Video) Validate() error {
	if _, err := govalidator.ValidateStruct(ref); err != nil {
		return err
	}
	return nil
}

type VideoRepository interface {
	Insert(video *Video) (*Video, error)
	Find(id string) (*Video, error)
}
