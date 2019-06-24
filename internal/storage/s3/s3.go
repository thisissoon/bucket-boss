package s3

import (
	"fmt"
	"strings"

	"github.com/thisissoon/bucket-boss/internal/config"
	"github.com/thisissoon/bucket-boss/internal/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

// Delete deletes a single object from the bucket
func (s *S3) Delete(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return fmt.Errorf("error delete object: %v", err)
		}
		return nil
	}
	return nil
}

func convertKeys(ks []string) []*s3.ObjectIdentifier {
	objects := []*s3.ObjectIdentifier{}
	for _, k := range ks {
		objects = append(objects, &s3.ObjectIdentifier{
			Key: aws.String(k),
		})
	}
	return objects
}

// DeleteMulti deletes multiple objects from the bucket.
// Large numbers of keys are batched into requests of a 1000 keys
func (s *S3) DeleteMulti(keys []string) error {
	for _, batch := range storage.Batcher(keys, 1000) {
		objects := convertKeys(batch)
		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &s3.Delete{
				Objects: objects,
				Quiet:   aws.Bool(false),
			},
		}
		_, err := s.svc.DeleteObjects(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				return fmt.Errorf("error deleting objects: %v", err)
			}
		}
	}
	return nil
}
