package svc

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_ecommerce/config"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
)

type ImageSvc interface {
	UploadImage(multipart.File, *multipart.FileHeader) (string, error)
}

type imageSvc struct {
	s3Client s3.Client
}

func NewImageSvc(cfg aws.Config) ImageSvc {
	s3Client := s3.NewFromConfig(cfg)

	return &imageSvc{
		s3Client: *s3Client,
	}
}

func (i *imageSvc) UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	bucketName := config.GetString("S3_BUCKET_NAME")
	uuid := uuid.NewString()
	splitted := strings.Split(fileHeader.Filename, ".")
	ext := splitted[len(splitted)-1]
	fileName := fmt.Sprintf("%s.%s", uuid, ext)

	if fileHeader.Size < 10*1024 || fileHeader.Size > 2*1024*1024 {
		return "", customErr.NewBadRequestError("file size must be between 10KB and 2MB")
	}
	if fileHeader.Header["Content-Type"][0] != "image/jpeg" {
		return "", customErr.NewBadRequestError("file must be in JPEG format")
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	}

	_, err := i.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.ap-southeast-1.amazonaws.com/%s", bucketName, fileName)

	return url, nil
}
