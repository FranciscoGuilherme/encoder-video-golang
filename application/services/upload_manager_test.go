package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"os"
	"testing"
	"time"
	"uuid"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type UploadManagerTestSuite struct {
	suite.Suite
	video *domain.Video
	repo  repositories.VideoRepository
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestUploadManagerSuite(t *testing.T) {
	suite.Run(t, new(UploadManagerTestSuite))
}

func (s *UploadManagerTestSuite) SetupTest() {
	db := database.NewDbTest(false)
	defer db.Close()

	s.video = domain.NewVideo()
	s.video.ID = uuid.NewV4().String()
	s.video.FilePath = "terra-test.mp4"
	s.video.CreatedAt = time.Now()
	s.repo = repositories.NewVideoRepository(db)
}

func (s *UploadManagerTestSuite) TestVideoServiceUpload() {
	videoService := services.NewVideoService()
	videoService.Video = s.video
	videoService.VideoRepository = s.repo

	err := videoService.Download("codeeducationtest")
	s.Require().Nil(err)

	err = videoService.Fragment()
	s.Require().Nil(err)

	err = videoService.Encode()
	s.Require().Nil(err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "codeeducationtest"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + s.video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload
	s.Require().Equal(result, "upload completed")

	err = videoService.Finish()
	s.Require().Nil(err)
}
