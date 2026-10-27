package domain

import (
	"time"
	"uuid"

	"github.com/asaskevich/govalidator"
)

type JobStatus string

type AvailableStatus struct {
	Starting    JobStatus
	Downloading JobStatus
	Fragmenting JobStatus
	Encoding    JobStatus
	Finishing   JobStatus
	Completed   JobStatus
	Uploading   JobStatus
	Failed      JobStatus
}

var JobStatusList = AvailableStatus{
	Starting:    JobStatus("STARTING"),
	Downloading: JobStatus("DOWNLOADING"),
	Fragmenting: JobStatus("FRAGMENTING"),
	Encoding:    JobStatus("ENCODING"),
	Finishing:   JobStatus("FINISHING"),
	Completed:   JobStatus("COMPLETED"),
	Uploading:   JobStatus("UPLOADING"),
	Failed:      JobStatus("FAILED"),
}

type Job struct {
	ID               string    `json:"job_id"             valid:"uuid"      gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status"             valid:"notnull"`
	Video            *Video    `json:"video"              valid:"-"`
	VideoID          string    `json:"-"                  valid:"-"         gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `json:"-"                  valid:"-"`
	CreatedAt        time.Time `json:"created_at"         valid:"-"`
	UpdatedAt        time.Time `json:"updated_at"         valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewJob(output string, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status:           string(JobStatusList.Starting),
		Video:            video,
	}

	job.prepare()

	err := job.Validate()
	if err != nil {
		return &Job{}, err
	}

	return &job, nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)
	if err != nil {
		return err
	}

	return nil
}
