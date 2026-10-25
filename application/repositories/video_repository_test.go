package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"
	"uuid"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type VideoRepositoryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestVideoRepositorySuite(t *testing.T) {
	suite.Run(t, new(VideoRepositoryTestSuite))
}

func (s *VideoRepositoryTestSuite) SetupTest() {
	s.db = database.NewDbTest(false)
}

func (s *VideoRepositoryTestSuite) TestVideoRepositoryDbInsert() {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(s.db)
	repo.Insert(video)
	v, err := repo.Find(video.ID)

	s.Require().NotEmpty(v.ID)
	s.Require().Nil(err)
	s.Require().Equal(v.ID, video.ID)
}

func (s *VideoRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}
