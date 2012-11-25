// Simple XRD parser
//
// XRD spec: http://docs.oasis-open.org/xri/xrd/v1.0/xrd-1.0.html
package xrd

import (
	"encoding/xml"
)

type XRD struct {
	Subject  string
	Expires  string
	Alias    []string
	Link     []Link
	Property []Property
}

type Property struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type Link struct {
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Href     string `xml:"href,attr"`
	Title    []Title
	Property []Property
	Template string `xml:"template,attr"`
}

// Parse the XRD using xml.Unmarshal
func ParseXRD(blob []byte) (*XRD, error) {
	parsed := XRD{}
	err := xml.Unmarshal(blob, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Return the first *Link with rel=rel
func (self *XRD) GetLinkByRel(rel string) *Link {
	for _, link := range self.Link {
		if link.Rel == rel {
			return &link
		}
	}
	return nil
}
