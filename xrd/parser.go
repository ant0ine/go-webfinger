// Package xrd provides a simple XRD parser.
//
// XRD spec: http://docs.oasis-open.org/xri/xrd/v1.0/xrd-1.0.html
package xrd

import (
	"encoding/xml"
)

// XRD is an Extensible Resource Descriptor, specifying properties and related
// links for a resource.
type XRD struct {
	Subject  string
	Expires  string
	Alias    []string
	Link     []Link
	Property []Property
}

// Property is a property of a resource.
type Property struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Title is a human-readable description of a Link.
type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Link is a link to a related resource.
type Link struct {
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Href     string `xml:"href,attr"`
	Title    []Title
	Property []Property
	Template string `xml:"template,attr"`
}

// ParseXRD parses the XRD using xml.Unmarshal.
func ParseXRD(blob []byte) (*XRD, error) {
	parsed := XRD{}
	err := xml.Unmarshal(blob, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// GetLinkByRel returns the first *Link with the specified rel value.
func (self *XRD) GetLinkByRel(rel string) *Link {
	for _, link := range self.Link {
		if link.Rel == rel {
			return &link
		}
	}
	return nil
}
