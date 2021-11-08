package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutPutBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video"  valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:notnull"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdateAt         time.Time `json:"update_at" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := Job{
		OutPutBucketPath: output,
		Status:           status,
		Video:            video,
	}

	job.prepare()
	err := job.Validate()
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdateAt = time.Now()
}

func (job *Job) Validate() error {

	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil

}
