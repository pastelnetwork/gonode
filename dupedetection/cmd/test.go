package main

import (
	"context"
	"flag"
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

	_, err = client.ImageRarenessScore(context.Background(), data, "png", "testBlockHash", "240669", "2022-03-31 16:55:28", "testPastelID", "jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis", "jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F", "jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ", false, "")
	if err != nil {
		panic("do ImageRarenessScore() error: " + err.Error())
	}
}
