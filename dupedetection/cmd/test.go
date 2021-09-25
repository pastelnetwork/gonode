package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
)

func main() {
	config := &ddclient.Config{
		Host:       "localhost",
		Port:       50052,
		DDFilesDir: "/tmp",
	}
	client := ddclient.NewDDServerClient(config)
	data, err := ioutil.ReadFile("test.png")
	if err != nil {
		panic("read file error" + err.Error())
	}

	_, err = client.ImageRarenessScore(context.Background(), data, "png")
	if err != nil {
		panic("do ImageRarenessScore() error" + err.Error())
	}

	fmt.Println("Successful")

}
