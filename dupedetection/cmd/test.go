package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
)

func main() {
	var (
		filePath = flag.String("path", "test.png", "path to image")
	)
	//flag.Usage = usage
	flag.Parse()
	config := &ddclient.Config{
		Host:       "localhost",
		Port:       50052,
		DDFilesDir: "/tmp",
	}
	client := ddclient.NewDDServerClient(config)
	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic("read file error" + err.Error())
	}

	res, err := client.ImageRarenessScore(context.Background(), data, "png", "testBlockHash", "testPastelID")
	if err != nil {
		panic("do ImageRarenessScore() error" + err.Error())
	}

	dataJSON, _ := json.Marshal(res)

	fmt.Println(string(dataJSON))

}
