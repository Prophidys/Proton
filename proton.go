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
	Auth() bool
	ListBucket() string
	CreateBucket() string
	DeleteBucket() string
	Put() string
	Get() string
	Del() string
}

type ObjectStoreService struct {
	ObjectStore string
	AuthInfo    json.RawMessage
}

type AWS_S3 struct {
	AccessKey       string
	SecretAccessKey string
}

type OS_Swift struct {
	Login       string
	Password    string
	TenantID    string
	EndPointURL string
}

func main() {
	loadConfig()
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
				listBuckets()
			},
		},
		{
			Name:      "delete",
			ShortName: "d",
			Usage:     "delete a bucket/container",
			Action: func(c *cli.Context) {
				deleteBucket(c.Args().First())
			},
		},
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "create a bucket/container",
			Action: func(c *cli.Context) {
				createBucket(c.Args().First())
			},
		},
	}

	app.Run(os.Args)
}

func loadConfig() {
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
	var dst interface{}
	switch oss.ObjectStore {
	case "s3":
		dst = new(AWS_S3)
	case "swift":
		dst = new(OS_Swift)
	}

	if err != nil {
		fmt.Println("error:", err)
	}
	json.Unmarshal(oss.AuthInfo, dst)
	if oss.ObjectStore == "s3" {
		var dstf *AWS_S3 = dst.(*AWS_S3)
		fmt.Println(fmt.Sprintf("%+v %+v", dstf.AccessKey, dstf.SecretAccessKey))
		os.Setenv("AWS_ACCESS_KEY", dstf.AccessKey)
		os.Setenv("AWS_SECRET_ACCESS_KEY", dstf.SecretAccessKey)
	}
}

func listBuckets() {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.USEast)
	resp, err := client.ListBuckets()

	if err != nil {
		log.Fatal(err)
	}

	for i := range resp.Buckets {
		fmt.Println(fmt.Sprintf("%+v", resp.Buckets[i].Name))
	}
}

func createBucket(bucketName string) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.USEast)
	test := client.Bucket(bucketName)
	fmt.Println(fmt.Sprintf("%+v", test))
	err2 := test.PutBucket(s3.Private)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func deleteBucket(bucketName string) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.USEast)
	test := client.Bucket(bucketName)
	fmt.Println(fmt.Sprintf("%+v", test))
	err2 := test.DelBucket()
	if err2 != nil {
		log.Fatal(err2)
	}
}
