package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type JobTestSuite struct {
	suite.Suite
	video *domain.Video
}

func TestJobSuite(t *testing.T) {
	suite.Run(t, new(JobTestSuite))
}

func (s *JobTestSuite) SetupTest() {
	s.video = domain.NewVideo()
	s.video.ID = uuid.NewV4().String()
	s.video.FilePath = "path"
	s.video.CreatedAt = time.Now()
}

func (s *JobTestSuite) TestNewJob() {
	job, err := domain.NewJob("path", "Converted", s.video)
	s.Require().NotNil(job)
	s.Require().Nil(err)
}
