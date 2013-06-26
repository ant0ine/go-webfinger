package xrd

import (
	"github.com/ant0ine/go-webfinger/jrd"
)

// ConvertToJRD converts the XRD to JRD.
func (self *XRD) ConvertToJRD() *jrd.JRD {
	obj := jrd.JRD{}

	obj.Subject = self.Subject
	obj.Expires = self.Expires
	obj.Aliases = self.Alias

	for _, link := range self.Link {
		obj.Links = append(obj.Links, convertLink(&link))
	}

	obj.Properties = make(map[string]interface{})
	for _, prop := range self.Property {
		obj.Properties[prop.Type] = prop.Value
	}

	return &obj
}

func convertLink(link *Link) jrd.Link {
	obj := jrd.Link{}

	obj.Rel = link.Rel
	obj.Type = link.Type
	obj.Href = link.Href

	obj.Titles = make(map[string]string)
	for _, title := range link.Title {
		key := title.Lang
		if key == "" {
			key = "default"
		}
		obj.Titles[key] = title.Value
	}

	obj.Properties = make(map[string]interface{})
	for _, prop := range link.Property {
		obj.Properties[prop.Type] = prop.Value
	}
	obj.Template = link.Template

	return obj
}
