package main

import (
	"fmt"
	"github.com/ant0ine/go-webfinger"
	"os"
)

func main() {
	email := os.Args[1]
	user_jrd, err := webfinger.GetUserJRD(email)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("User JRD: %+v", user_jrd)
}
