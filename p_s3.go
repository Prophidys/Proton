package main

import (
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
	"os"
)

type AWS_S3 struct {
	AccessKey         string
	SecretAccessKey   string
	AuthHandler       aws.Auth
	ConnectionHandler *s3.S3
}

func (s *AWS_S3) Auth() string {
	var err error
	os.Setenv("AWS_ACCESS_KEY", s.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", s.SecretAccessKey)
	s.AuthHandler, err = aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	s.ConnectionHandler = s3.New(s.AuthHandler, aws.EUWest)
	return "S3 Auth"
}

func (s *AWS_S3) ListBuckets() string {
	resp, err := s.ConnectionHandler.ListBuckets()
	if err != nil {
		log.Fatal(err)
	}
	for i := range resp.Buckets {
		fmt.Println(fmt.Sprintf("%+v", resp.Buckets[i].Name))
	}
	return "S3 ListBucket"
}

func (s *AWS_S3) CreateBucket(bucketName string) string {
	bucket := s.ConnectionHandler.Bucket(bucketName)
	err := bucket.PutBucket(s3.Private)
	if err != nil {
		log.Fatal(err)
	}
	return "S3 Create Bucket"
}

func (s *AWS_S3) DeleteBucket(bucketName string) string {
	bucket := s.ConnectionHandler.Bucket(bucketName)
	err := bucket.DelBucket()
	if err != nil {
		log.Fatal(err)
	}
	return "S3 Delete Bucket"
}

func (s *AWS_S3) Put(src string, dst string) string {
	return "S3 Put Object"
}

func (s *AWS_S3) Get() string {
	return "S3 Get Object"
}

func (s *AWS_S3) Del() string {
	return "S3 Put Object"
}
