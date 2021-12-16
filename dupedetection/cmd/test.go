package main

import (
	"context"
	"encoding/json"
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

	res, err := client.ImageRarenessScore(context.Background(), data, "png", "testBlockHash", "testPastelID")
	if err != nil {
		panic("do ImageRarenessScore() error" + err.Error())
	}

	dataJson, _ := json.Marshal(res)

	fmt.Println(string(dataJson))

}
