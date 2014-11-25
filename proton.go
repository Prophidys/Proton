package main

import "fmt"
import "github.com/mitchellh/goamz/aws"
import "encoding/json"
import "os"
import "github.com/codegangsta/cli"
import "io/ioutil"

// import "github.com/ncw/swift"
import "github.com/mitchellh/goamz/s3"
import "log"

type ObjectStore interface {
	Auth() string
	ListBuckets() string
	CreateBucket(bucketName string) string
	DeleteBucket(bucketName string) string
	Put() string
	Get() string
	Del() string
}

type ObjectStoreService struct {
	ObjectStore string
	AuthInfo    json.RawMessage
}

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

func (s *AWS_S3) Put() string {
	return "S3 Put Object"
}

func (s *AWS_S3) Get() string {
	return "S3 Get Object"
}

func (s *AWS_S3) Del() string {
	return "S3 Put Object"
}

type OS_Swift struct {
	Login       string
	Password    string
	TenantID    string
	EndPointURL string
}

func main() {
	objstore := loadConfig()
	objstore.Auth()
	app := cli.NewApp()
	app.Name = "proton"
	app.Usage = "Object Storage Abstractor"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list all avaiable buckets/container",
			Action: func(c *cli.Context) {
				objstore.ListBuckets()
			},
		},
		{
			Name:      "delete",
			ShortName: "d",
			Usage:     "delete a bucket/container",
			Action: func(c *cli.Context) {
				objstore.DeleteBucket(c.Args().First())
			},
		},
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "create a bucket/container",
			Action: func(c *cli.Context) {
				objstore.CreateBucket(c.Args().First())
			},
		},
	}

	app.Run(os.Args)
}

func loadConfig() ObjectStore {
	var oss ObjectStoreService
	var filename = ""
	if os.Getenv("PROTON_CONFIG") != "" {
		filename = os.Getenv("PROTON_CONFIG")
	} else if _, err := os.Stat("/etc/proton/proton.cfg"); !os.IsNotExist(err) {
		filename = "/etc/proton/proton.cfg"
	} else if _, err := os.Stat("~/.proton/proton.cfg"); !os.IsNotExist(err) {
		filename = "~/.proton/proton.cfg"
	} else if _, err := os.Stat("proton.cfg"); !os.IsNotExist(err) {
		filename = "proton.cfg"
	} else {
		log.Fatal("Config file not found")
	}
	file, _ := ioutil.ReadFile(filename)
	err := json.Unmarshal(file, &oss)
	var objstore ObjectStore
	switch oss.ObjectStore {
	case "s3":
		objstore = new(AWS_S3)
	case "swift":
		//		objstore = new(OS_Swift)
	}

	if err != nil {
		fmt.Println("error:", err)
	}
	json.Unmarshal(oss.AuthInfo, objstore)
	return objstore
}
