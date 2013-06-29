
Go-Webfinger
============

*Go client for the Webfinger protocol*

[![Build Status](https://travis-ci.org/ant0ine/go-webfinger.png?branch=master)](https://travis-ci.org/ant0ine/go-webfinger)

**Go-Webfinger** is a Go client for the Webfinger protocol.

*It is a work in progress, the API is not frozen.
We're trying to catchup with the last draft of the protocol:
http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-14
and to support the http://webfist.org *

Install
-------

This package is "go-gettable", just do:

    go get github.com/ant0ine/go-webfinger

Example
-------

    package main

    import (
            "fmt"
            "github.com/ant0ine/go-webfinger"
            "os"
    )

    func main() {
            email := os.Args[1]

            client := webfinger.NewClient(nil)

            resource, err := webfinger.MakeResource(email)
            if err != nil {
                    panic(err)
            }

            jrd, err := client.GetJRD(resource)
            if err != nil {
                    fmt.Println(err)
                    return
            }

            fmt.Printf("JRD: %+v", jrd)
    }

Documentation
-------------

- [Online Documentation (godoc.org)](http://godoc.org/github.com/ant0ine/go-webfinger)

Author
------
- [Antoine Imbert](https://github.com/ant0ine)

Contributors
------------

- Thanks [Will Norris](https://github.com/willnorris) for the major update to support draft-14, and the GAE compat!


[MIT License](https://github.com/ant0ine/go-webfinger/blob/master/LICENSE)
