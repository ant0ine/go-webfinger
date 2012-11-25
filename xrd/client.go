package xrd

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// XRD spec: http://docs.oasis-open.org/xri/xrd/v1.0/xrd-1.0.html

type Link struct {
	Rel      string `xml:"rel,attr"`
	Template string `xml:"template,attr"`
	Type     string `xml:"type,attr"`
	Href     string `xml:"href,attr"`
}

type XRD struct {
	Subject string
	Alias   string
	Expires string
	Link    []Link
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

	parsed := XRD{}
	err = xml.Unmarshal(content, &parsed)
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
