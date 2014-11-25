package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Local struct {
	Path string
}

func (l *Local) Auth() string {
	return "Auth Local"
}

func (l *Local) ListBuckets() string {
	listdir, _ := ioutil.ReadDir(l.Path)
	for i := range listdir {
		fmt.Println(listdir[i])
	}
	return "Local List Buckets"
}

func (l *Local) CreateBucket(bucketName string) string {
	os.MkdirAll(l.Path+"/"+bucketName, 0777)
	return "Local Create Bucket"
}

func (l *Local) DeleteBucket(bucketName string) string {
	return "Local Delete Bucket"
}

func (l *Local) Put(src string, dst string) string {
	return "Local Put Object"
}

func (l *Local) Get() string {
	return "Local Get Object"
}

func (l *Local) Del() string {
	return "Local Put Object"
}
