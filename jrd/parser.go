// Simple JRD Parser
//
// JRD spec: http://tools.ietf.org/html/rfc6415#appendix-A
package jrd

import (
	"encoding/json"
)

type JRD struct {
	Subject    string
	Expires    string
	Aliases    []string
	Links      []Link
	Properties map[string]interface{}
}

type Link struct {
	Rel        string
	Type       string
	Href       string
	Titles     map[string]string
	Properties map[string]interface{}
	Template   string
}

// Parse the JRD using json.Unmarshal
func ParseJRD(blob []byte) (*JRD, error) {
	jrd := JRD{}
	err := json.Unmarshal(blob, &jrd)
	if err != nil {
		return nil, err
	}
	return &jrd, nil
}

// Return the first *Link with rel=rel
func (self *JRD) GetLinkByRel(rel string) *Link {
	for _, link := range self.Links {
		if link.Rel == rel {
			return &link
		}
	}
	return nil
}
