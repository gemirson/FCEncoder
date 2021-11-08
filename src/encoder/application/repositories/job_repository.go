package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDd struct {
	Db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *VideoRepositoryDd {
	return &VideoRepositoryDd{
		Db: db,
	}
}

func (repo JobRepositoryDd) Insert(job *domain.Job) (*domain.Job, error) {

	err := repo.Db.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repo JobRepositoryDd) Find(id string) (*domain.Job, error) {

	var job domain.Job
	repo.Db.Preload("Video").First(&job, "id=?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("this video not found ")
	}

	return &job, nil

}

func (repo JobRepositoryDd) Update(job *domain.Job) (*domain.Job, error) {

	err := repo.Db.Save(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
