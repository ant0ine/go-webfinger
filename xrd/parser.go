package xrd

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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

func GetXRD(url string) (*XRD, error) {
	// TODO follow redirect
	// TODO try http if https fails
	// TODO verify signature if not https
	// TODO extract http cache info

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	parsed, err := ParseXRD(content)
	if err != nil {
		return nil, err
	}

	return parsed, nil
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

