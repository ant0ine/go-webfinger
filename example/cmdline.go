package main

import (
	"fmt"
	"github.com/ant0ine/go-webfinger"
	"os"
)

func main() {
	email := os.Args[1]
	resource, err := webfinger.MakeResource(email)
	if err != nil {
		panic(err)
	}
	jrd, err := resource.GetJRD()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("JRD: %+v", jrd)
}
