package main

import (
	"fmt"
	"github.com/ant0ine/go-webfinger"
	"os"
)

func main() {
	email := os.Args[1]
	user_xrd, err := webfinger.GetUserXRD(email)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("User XRD: %+v", user_xrd)
}
