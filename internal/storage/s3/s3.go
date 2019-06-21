package s3

import (
	"fmt"
	"github.com/thisissoon/bucket-boss/internal/config"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	svc    *s3.S3
	bucket string
}

func NewS3(cfg config.AWS) (*S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating an AWS session: %v", err)
	}
	return &S3{
		svc:    s3.New(session),
		bucket: cfg.BucketName,
	}, nil
}

// List returns s3 object keys with a given extension
// returns all files if the extension is empty string
func (s *S3) List(ext string) ([]string, error) {
	var keys []string
	err := s.svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &s.bucket,
	}, func(p *s3.ListObjectsOutput, last bool) bool {
		for _, obj := range p.Contents {
			if ext == "" {
				keys = append(keys, *obj.Key)
			} else {
				if strings.HasSuffix(*obj.Key, "."+ext) {
					keys = append(keys, *obj.Key)
				}
			}
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}
	return keys, nil
}
