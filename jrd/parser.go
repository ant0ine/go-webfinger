// Package jrd provides a simple JRD parser.
//
// Following this JRD spec: http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-14#section-4.4
//
package jrd

import (
	"encoding/json"
)

// JRD is a JSON Resource Descriptor, specifying properties and related links
// for a resource.
type JRD struct {
	Subject    string
	Aliases    []string
	Properties map[string]interface{}
	Links      []Link
}

// Link is a link to a related resource.
type Link struct {
	Rel        string
	Type       string
	Href       string
	Titles     map[string]string
	Properties map[string]interface{}
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

// GetProperty Returns the property value as a string.
// Per spec a property value can be null, empty string is returned in this case.
func (self *JRD) GetProperty(uri string) string {
	if self.Properties[uri] == nil {
		return ""
	}
	return self.Properties[uri].(string)
}

// GetProperty Returns the property value as a string.
// Per spec a property value can be null, empty string is returned in this case.
func (self *Link) GetProperty(uri string) string {
	if self.Properties[uri] == nil {
		return ""
	}
	return self.Properties[uri].(string)
}
