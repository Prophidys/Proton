package main

import "fmt"
import "github.com/mitchellh/goamz/aws"
import "encoding/json"
import "os"
import "github.com/codegangsta/cli"

// import "github.com/ncw/swift"
import "github.com/mitchellh/goamz/s3"
import "log"

type Configuration struct {
	ObjectEndPoint  string
	AccessKey       string
	SecretAccessKey string
	Region          string
}

func main() {
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

func listBuckets() {
	file, _ := os.Open("proton.cfg")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("error:", err)
	}

	os.Setenv("AWS_ACCESS_KEY", configuration.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", configuration.SecretAccessKey)

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
	file, _ := os.Open("proton.cfg")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("error:", err)
	}

	os.Setenv("AWS_ACCESS_KEY", configuration.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", configuration.SecretAccessKey)

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
	file, _ := os.Open("proton.cfg")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("error:", err)
	}

	os.Setenv("AWS_ACCESS_KEY", configuration.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", configuration.SecretAccessKey)

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