package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"time"
)

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3(session *session.Session, timeout time.Duration) *S3 {
	return &S3{
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

func (ref S3) Create(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, ref.timeout)
	defer cancel()

	if _, err := ref.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := ref.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (ref S3) UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ref.timeout)
	defer cancel()

	res, err := ref.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	return res.Location, nil
}

func (ref S3) DownloadObject(ctx context.Context, bucket, filename string, body io.WriterAt) error {
	if _, err := ref.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}

func (ref S3) DeleteObject(ctx context.Context, bucket, filename string) error {
	if _, err := ref.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if err := ref.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (ref S3) ListObjects(ctx context.Context, bucket string) ([]*Object, error) {
	ctx, cancel := context.WithTimeout(ctx, ref.timeout)
	defer cancel()

	res, err := ref.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	objects := make([]*Object, len(res.Contents))

	for i, obj := range res.Contents {
		objects[i] = &Object{
			Key:        *obj.Key,
			Size:       *obj.Size,
			ModifiedAt: *obj.LastModified,
		}
	}
	return objects, nil
}

func (ref S3) FetchObject(ctx context.Context, bucket, filename string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, ref.timeout)
	defer cancel()

	res, err := ref.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
