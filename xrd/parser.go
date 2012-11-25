package xrd

import (
	"encoding/xml"
)

// XRD spec: http://docs.oasis-open.org/xri/xrd/v1.0/xrd-1.0.html

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

func ParseXRD(blob []byte) (*XRD, error) {
	parsed := XRD{}
	err := xml.Unmarshal(blob, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (self *XRD) GetLinkByRel(rel string) *Link {
	for _, link := range self.Link {
		if link.Rel == rel {
			return &link
		}
	}
	return nil
}

func (self *XRD) LrddTemplate() string {
	link := self.GetLinkByRel("lrdd")
	if link == nil {
		return ""
	}
	return link.Template
}
