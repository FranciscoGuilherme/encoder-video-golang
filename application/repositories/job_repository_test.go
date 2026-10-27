package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type JobRepositoryTestSuite struct {
	suite.Suite
	db    *gorm.DB
	video *domain.Video
}

func TestJobRepositorySuite(t *testing.T) {
	suite.Run(t, new(JobRepositoryTestSuite))
}

func (s *JobRepositoryTestSuite) SetupTest() {
	s.db = database.NewDbTest(false)
	s.video = domain.NewVideo()
	s.video.ID = uuid.NewV4().String()
	s.video.FilePath = "path"
	s.video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: s.db}
	repo.Insert(s.video)
}

func (s *JobRepositoryTestSuite) TestJobRepositoryDbInsert() {
	job, err := domain.NewJob("output_path", s.video)
	s.Require().Nil(err)

	repoJob := repositories.JobRepositoryDb{Db: s.db}
	repoJob.Insert(job)
	j, err := repoJob.Find(job.ID)

	s.Require().NotEmpty(j.ID)
	s.Require().Nil(err)
	s.Require().Equal(j.ID, job.ID)
	s.Require().Equal(j.VideoID, s.video.ID)
}

func (s *JobRepositoryTestSuite) TestJobRepositoryDbUpdate() {
	job, err := domain.NewJob("output_path", s.video)
	s.Require().Nil(err)

	repoJob := repositories.JobRepositoryDb{Db: s.db}
	repoJob.Insert(job)
	job.Status = string(domain.JobStatusList.Completed)
	repoJob.Update(job)
	j, err := repoJob.Find(job.ID)

	s.Require().NotEmpty(j.ID)
	s.Require().Nil(err)
	s.Require().Equal(j.Status, job.Status)
}

func (s *JobRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}
