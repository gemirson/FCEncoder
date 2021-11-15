package services

import (
	"encoder/domain"
	"encoder/framework/helpers"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job domain.Job, workerID int) {

	for message := range messageChannel {

		err := helpers.IsJson(string(message.Body))

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoServiceAzure.Video)
		jobService.VideoServiceAzure.Video.ID = uuid.NewV4().String()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoServiceAzure.Video.Validate()
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoServiceAzure.InsertVideo()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.createJobWithStatusSTARTING()

		err = jobService.InsertJob()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.Start()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)

	}

}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
	return result
}
