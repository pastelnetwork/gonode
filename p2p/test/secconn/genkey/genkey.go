package main

import (
	"encoding/base64"
	"fmt"

	"github.com/otrv4/ed448"
)

// Create a preconfig keys, used in common package
func main() {
	curve := ed448.NewCurve()
	pri, pub, _ := curve.GenerateKeys()

	fmt.Println(base64.StdEncoding.EncodeToString(pri[:]))
	fmt.Println(base64.StdEncoding.EncodeToString(pub[:]))
}
