package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job               *domain.Job
	JobRepository     repositories.JobRepository
	VideoServiceAzure VideoServiceAzure
}

func NewJobService() JobService {
	return JobService{}
}

func (j *JobService) Start() error {

	err := j.changeJobStatus("DOWNLOADING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoServiceAzure.Download(os.Getenv("inputBucketName"))

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FRAGMENTING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoServiceAzure.Fragment()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("ENCODING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoServiceAzure.Encode()

	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FINISHING")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoServiceAzure.Finish()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("COMPLETED")

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {

	err := j.changeJobStatus("UPLOADING")
	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUploadAzure()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = j.VideoServiceAzure.ExtractedPathSource()

	concurrency, _ := strconv.Atoi("CONCURRENCY_UPLOAD")
	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(concurrency, doneUpload, j.VideoServiceAzure.AzureContainer)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err

}

func (j *JobService) changeJobStatus(status string) error {
	var err error
	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {

	j.Job.Status = "FAILED"
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return error
}
