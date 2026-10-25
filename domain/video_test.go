package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type VideoTestSuite struct {
	suite.Suite
	video *domain.Video
}

func TestVideoSuite(t *testing.T) {
	suite.Run(t, new(VideoTestSuite))
}

func (s *VideoTestSuite) SetupTest() {
	s.video = domain.NewVideo()
}

func (s *VideoTestSuite) TestValidateIfVideoIsEmpty() {
	err := s.video.Validate()
	s.Require().Error(err)
}

func (s *VideoTestSuite) TestVideoIdIsNotAUuid() {
	s.video.ID = "abc"
	s.video.ResourceID = "a"
	s.video.FilePath = "path"
	s.video.CreatedAt = time.Now()

	err := s.video.Validate()
	s.Require().Error(err)
}

func (s *VideoTestSuite) TestVideoValidation() {
	s.video.ID = uuid.NewV4().String()
	s.video.ResourceID = "a"
	s.video.FilePath = "path"
	s.video.CreatedAt = time.Now()

	err := s.video.Validate()
	s.Require().Nil(err)
}
