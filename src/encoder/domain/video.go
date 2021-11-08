package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoder_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceId string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path"   valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	Jobs       []*Job    `json:"-" valid:"-" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

}

func NewVideo() *Video {
	return &Video{}

}

func CreateInstanceVideo(ID string, ResourceId string, FilePath string, CreatedAt time.Time) *Video {
	return &Video{
		ID:         ID,
		ResourceId: ResourceId,
		FilePath:   FilePath,
		CreatedAt:  CreatedAt,
	}
}

func (video *Video) Validate() error {

	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil

}
