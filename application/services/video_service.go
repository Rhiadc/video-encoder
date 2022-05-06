package services

import (
	"context"
	"fmt"
	"github.com/video-encoder/domain"
	MyS3 "github.com/video-encoder/pkg/s3"
	"log"
	"os"
)

type VideoService struct {
	VideoRepository domain.VideoRepository
	MyS3            *MyS3.S3
}

func NewVideoService(videoRepository domain.VideoRepository, MyS3 *MyS3.S3) VideoService {
	return VideoService{VideoRepository: videoRepository, MyS3: MyS3}
}

func (ref VideoService) Create(ctx context.Context, bucket string) {
	if err := ref.MyS3.Create(ctx, bucket); err != nil {
		log.Fatalln(err)
	}
	//logger info "created"
}

func (ref VideoService) Download(ctx context.Context, bucket string, video domain.Video) error {
	fileName := fmt.Sprintf("%s/%s.mp4", os.Getenv("LOCAL_STORAGE_PATH"), video.ID)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := ref.MyS3.DownloadObject(ctx, bucket, fileName, file); err != nil {
		return err
	}
	return nil
}

func (ref VideoService) Upload(ctx context.Context, bucket string) error {
	file, err := os.Open("./assets/id.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	url, err := ref.MyS3.UploadObject(ctx, bucket, "id.txt", file)
	if err != nil {
		return err
	}
	log.Println("upload object:", url)

	return nil
}
