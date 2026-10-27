package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func NewJobService(job *domain.Job, repo repositories.JobRepository, service VideoService) JobService {
	return JobService{
		Job:           job,
		JobRepository: repo,
		VideoService:  service,
	}
}

func (j *JobService) Start() error {
	err := j.changeJobStatus(domain.JobStatusList.Downloading)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("inputBucketName"))
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusList.Fragmenting)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusList.Encoding)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.upload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusList.Finishing)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusList.Completed)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) upload() error {
	var doneUpload chan string
	var uploadResult string

	err := j.changeJobStatus(domain.JobStatusList.Uploading)
	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))

	doneUpload = make(chan string)
	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult = <-doneUpload
	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err
}

func (j *JobService) changeJobStatus(status domain.JobStatus) error {
	var err error

	j.Job.Status = string(status)
	j.Job, err = j.JobRepository.Update(j.Job)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = string(domain.JobStatusList.Failed)
	j.Job.Error = error.Error()
	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return error
}
