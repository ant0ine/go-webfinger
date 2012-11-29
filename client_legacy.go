package webfinger

import (
	"errors"
	"github.com/ant0ine/go-webfinger/jrd"
	"log"
	"net/url"
	"strings"
)

func (self *Client) find_JRD(urls []string) (*jrd.JRD, error) {
	for _, try := range urls {
		try_obj, err := url.Parse(try)
		if err != nil {
			log.Print(err)
			continue
		}
		obj, err := self.fetch_JRD(try_obj)
		if err != nil {
			log.Print(err)
			continue
		}
		return obj, nil
	}
	return nil, errors.New("JRD not found")
}

// Build a serie well known host JRD URLs from the domain
func (self *Client) LegacyHostJRDURLs(domain string) []string {
	return []string{
		// first JRD implementation
		"https://" + domain + "/.well-known/host-meta.json",
		// orignal spec: https://code.google.com/p/webfinger/wiki/WebFingerProtocol
		"https://" + domain + "/.well-known/host-meta",
	}
}

// Given a domain, this method gets the host meta JRD data,
// and returns the LRDD resource JRD template URL.
// It tries all the urls returned by client.LegacyHostJRDURLs.
func (self *Client) LegacyGetResourceJRDTemplateURL(domain string) (string, error) {
	// TODO implement heavy HTTP cache around this

	urls := self.LegacyHostJRDURLs(domain)

	host_jrd, err := self.find_JRD(urls)
	if err != nil {
		return "", err
	}

	link := host_jrd.GetLinkByRel("lrdd")
	if link == nil {
		return "", errors.New("cannot find the LRDD link in the JRD data")
	}

	template := link.Template
	if template == "" {
		return "", errors.New("cannot find the template in the JRD data")
	}

	return template, nil
}

// Get the JRD data for this resource.
// Implement the original WebFinger API, ie: first fetch the Host metadata,
// find the LRDD link, fetch the resource data and convert the XRD in JRD if necessary.
func (self *Client) LegacyGetJRD(resource *Resource) (*jrd.JRD, error) {

	template, err := self.LegacyGetResourceJRDTemplateURL(resource.Domain)
	if err != nil {
		return nil, err
	}

	log.Printf("template: %s", template)

	jrd_url := strings.Replace(template, "{uri}", url.QueryEscape(resource.AsURIString()), 1)

	log.Printf("User JRD URL: %s", jrd_url)

	resource_jrd, err := self.find_JRD([]string{jrd_url})
	if err != nil {
		return nil, err
	}

	return resource_jrd, nil
}
