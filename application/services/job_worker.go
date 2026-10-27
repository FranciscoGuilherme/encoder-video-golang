package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"os"
	"sync"
	"uuid"

	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

type JobWorkerInitParameters struct {
	job            *domain.Job
	messageChannel chan amqp.Delivery
	returnChannel  chan JobWorkerResult
	jobService     JobService
	workerId       int
}

var Mutex = &sync.Mutex{}

func JobWorker(data JobWorkerInitParameters) {
	for message := range data.messageChannel {
		err := utils.IsJson(string(message.Body))
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		Mutex.Lock()
		err = json.Unmarshal(message.Body, &data.jobService.VideoService.Video)
		data.jobService.VideoService.Video.ID = uuid.NewV4().String()
		Mutex.Unlock()
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = data.jobService.VideoService.Video.Validate()
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		Mutex.Lock()
		err = data.jobService.VideoService.InsertVideo()
		Mutex.Unlock()
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		data.job, err = domain.NewJob(
			os.Getenv("outputBucketName"),
			data.jobService.VideoService.Video,
		)
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		Mutex.Lock()
		_, err = data.jobService.JobRepository.Insert(data.job)
		Mutex.Unlock()
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		data.jobService.Job = data.job
		err = data.jobService.Start()
		if err != nil {
			data.returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		data.returnChannel <- returnJobResult(*data.job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
}
