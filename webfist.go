package webfinger

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ant0ine/go-webfinger/jrd"
)

const (
	webFistDefaultServer = "webfist.org"
	webFistRel = "http://webfist.org/spec/rel"
)

func (c *Client) webfistLookup(resource *Resource) (*jrd.JRD, error) {
	jrdURL := resource.JRDURL(c.WebFistServer, nil)
	webfistJRD, err := c.fetchJRD(jrdURL)
	if err != nil {
		return nil, err
	}

	link := webfistJRD.GetLinkByRel(webFistRel)
	if link == nil {
		return nil, fmt.Errorf("No WebFist link")
	}

	u, err := url.Parse(link.Href)
	if err != nil {
		return nil, err
	}

	log.Printf("Found WebFist link: %s", u)
	return c.fetchJRD(u)
}
