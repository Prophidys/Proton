package main

import (
	"fmt"
	"github.com/ncw/swift"
	"log"
)

type OS_Swift struct {
	Login        string
	Password     string
	TenantID     string
	EndPointURL  string
	Tenant       string
	SwiftHandler swift.Connection
}

func (s *OS_Swift) Auth() string {
	s.SwiftHandler = swift.Connection{
		UserName: s.Login,
		ApiKey:   s.Password,
		AuthUrl:  s.EndPointURL,
		Tenant:   s.Tenant,
	}
	return "Swift Auth"
}

func (s *OS_Swift) ListBuckets() string {
	containers, err := s.SwiftHandler.ContainerNames(nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := range containers {
		fmt.Println(fmt.Sprintf("%+v", containers[i]))
	}
	return "Swift ListBucket"
}

func (s *OS_Swift) CreateBucket(bucketName string) string {
	err := s.SwiftHandler.ContainerCreate(bucketName, nil)
	if err != nil {
		log.Fatal(err)
	}
	return "Swift Create Bucket"
}

func (s *OS_Swift) DeleteBucket(bucketName string) string {
	err := s.SwiftHandler.ContainerDelete(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	return "Swift Delete Bucket"
}

func (s *OS_Swift) Put(src string, dst string) string {
	//_, err := s.SwiftHandler.ObjectPut(container string, objectName string, contents io.Reader, checkHash bool, Hash string, contentType string, h Headers)
	return "Swift Put Object"
}

func (s *OS_Swift) Get() string {
	return "Swift Get Object"
}

func (s *OS_Swift) Del() string {
	return "Swift Put Object"
}
