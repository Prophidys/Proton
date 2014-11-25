package main

import "fmt"
import "encoding/json"
import "os"
import "github.com/codegangsta/cli"
import "io/ioutil"
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
		objstore = new(OS_Swift)
	}

	if err != nil {
		fmt.Println("error:", err)
	}
	json.Unmarshal(oss.AuthInfo, objstore)
	return objstore
}
