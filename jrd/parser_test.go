package jrd

import (
	"testing"
)

func TestParseJRD(t *testing.T) {

	// Adapted from spec http://tools.ietf.org/html/rfc6415#appendix-A
	blob := `
        {
              "subject":"http://blog.example.com/article/id/314",

              "aliases":[
                "http://blog.example.com/cool_new_thing",
                "http://blog.example.com/steve/article/7"],

              "properties":{
                "http://blgx.example.net/ns/version":"1.3",
                "http://blgx.example.net/ns/ext":null
              },

              "links":[
                {
                  "rel":"author",
                  "type":"text/html",
                  "href":"http://blog.example.com/author/steve",
                  "titles":{
                    "default":"About the Author",
                    "en-us":"Author Information"
                  },
                  "properties":{
                    "http://example.com/role":"editor"
                  }
                },
                {
                  "rel":"author",
                  "href":"http://example.com/author/john",
                  "titles":{
                    "default":"The other author"
                  }
                },
                {
                  "rel":"copyright"
                }
              ]
            }
        `
	obj, err := ParseJRD([]byte(blob))
	if err != nil {
		t.Fatal(err)
	}
	if obj.Subject != "http://blog.example.com/article/id/314" {
		t.Error()
	}
	if obj.GetProperty("http://blgx.example.net/ns/version") != "1.3" {
		t.Error()
	}
	if obj.GetProperty("http://blgx.example.net/ns/ext") != "" {
		t.Error()
	}
	if obj.GetLinkByRel("copyright") == nil {
		t.Error()
	}
	if obj.GetLinkByRel("author").Titles["default"] != "About the Author" {
		t.Error()
	}
	if obj.GetLinkByRel("author").GetProperty("http://example.com/role") != "editor" {
		t.Error()
	}
}
