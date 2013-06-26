// Package jrd provides a simple JRD parser.
//
// JRD spec: http://tools.ietf.org/html/rfc6415#appendix-A
package jrd

import (
	"encoding/json"
)

// JRD is a JSON Resource Descriptor, specifying properties and related links
// for a resource.
type JRD struct {
	Subject    string
	Expires    string
	Aliases    []string
	Links      []Link
	Properties map[string]interface{}
}

// Link is a link to a related resource.
type Link struct {
	Rel        string
	Type       string
	Href       string
	Titles     map[string]string
	Properties map[string]interface{}
	Template   string
}

// ParseJRD parses the JRD using json.Unmarshal.
func ParseJRD(blob []byte) (*JRD, error) {
	jrd := JRD{}
	err := json.Unmarshal(blob, &jrd)
	if err != nil {
		return nil, err
	}
	return &jrd, nil
}

// GetLinkByRel returns the first *Link with the specified rel value.
func (self *JRD) GetLinkByRel(rel string) *Link {
	for _, link := range self.Links {
		if link.Rel == rel {
			return &link
		}
	}
	return nil
}
