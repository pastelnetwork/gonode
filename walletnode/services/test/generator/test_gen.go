//go:generate go run test_gen.go

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func parse(inputFile string, outputFile string, config interface{}) {

	t, err := template.ParseFiles(inputFile)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("create file: ", err)
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}

func genSupernodeClientTest() {
	type args struct {
		Service string
		Client  string
		Prefix  string
	}

	data := map[string]args{
		"../sense_node_client_test.go": {
			Service: "senseregister",
			Client:  "sense_register",
			Prefix:  "RegisterSense",
		},
		"../nftregistration_node_client_test.go": {
			Service: "nftregister",
			Client:  "artwork_register",
			Prefix:  "RegisterNft",
		},
	}
	for f, d := range data {
		parse("../templates/node_client_test.tmpl", f, d)
	}

}

func main() {
	genSupernodeClientTest()
}
