package helper

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/irnafitriani/music/entity"
	"github.com/thedevsaddam/govalidator"
)

func Validate(data entity.HasRules) url.Values {
	opts := govalidator.Options{
		Data:  data,
		Rules: data.Rules(),
	}
	v := govalidator.New(opts)
	e := v.ValidateStruct()

	return e
}

func UploadS3(filePath string, s3Session *session.Session, bucketName string) (*s3manager.UploadOutput, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, errors.New("file not found")
	}

	uploader := s3manager.NewUploader(s3Session)

	file, _ := os.Open(filePath)
	defer file.Close()
	str := strings.Split(filePath, "/")
	fileName := str[len(str)-1]

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		file.Close()
		return nil, err
	}

	return output, nil
}
