// TODO
// * improve JRD output
// * do stuff with the JRD
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-webfinger"
	"io/ioutil"
	"log"
	"os"
)

func printHelp() {
	fmt.Println("webfinger [-vh] <resource uri>")
	flag.PrintDefaults()
	fmt.Println("example: webfinger -v bob@example.com") // same Bob as in the draft
}

func main() {

	// cmd line flags
	verbose := flag.Bool("v", false, "print details about the resolution")
	help := flag.Bool("h", false, "display this message")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	email := flag.Arg(0)

	if email == "" {
		printHelp()
		os.Exit(1)
	}

	log.SetFlags(0)

	client := webfinger.NewClient(nil)
	client.AllowHTTP = true

	jrd, err := client.Lookup(email, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bytes, err := json.MarshalIndent(jrd, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bytes)

	os.Exit(0)
}
