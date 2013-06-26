package webfinger

import (
	"errors"
	"github.com/ant0ine/go-webfinger/jrd"
	"log"
	"net/url"
	"strings"
)

func (self *Client) findJRD(urls []string) (*jrd.JRD, error) {
	for _, try := range urls {
		tryObj, err := url.Parse(try)
		if err != nil {
			log.Print(err)
			continue
		}
		obj, err := self.fetchJRD(tryObj)
		if err != nil {
			log.Print(err)
			continue
		}
		return obj, nil
	}
	return nil, errors.New("JRD not found")
}

// LegacyHostJRDURLs builds a series of well known host JRD URLs from the domain.
func (self *Client) LegacyHostJRDURLs(domain string) []string {
	return []string{
		// first JRD implementation
		"https://" + domain + "/.well-known/host-meta.json",
		// orignal spec: https://code.google.com/p/webfinger/wiki/WebFingerProtocol
		"https://" + domain + "/.well-known/host-meta",
	}
}

// LegacyGetResourceJRDTemplateURL gets the host meta JRD data for the specified domain,
// and returns the LRDD resource JRD template URL.
// It tries all the urls returned by client.LegacyHostJRDURLs.
func (self *Client) LegacyGetResourceJRDTemplateURL(domain string) (string, error) {
	// TODO implement heavy HTTP cache around this

	urls := self.LegacyHostJRDURLs(domain)

	hostJRD, err := self.findJRD(urls)
	if err != nil {
		return "", err
	}

	link := hostJRD.GetLinkByRel("lrdd")
	if link == nil {
		return "", errors.New("cannot find the LRDD link in the JRD data")
	}

	template := link.Template
	if template == "" {
		return "", errors.New("cannot find the template in the JRD data")
	}

	return template, nil
}

// LegacyGetJRD gets the JRD data for this resource.
// Implement the original WebFinger API, ie: first fetch the Host metadata,
// find the LRDD link, fetch the resource data and convert the XRD in JRD if necessary.
func (self *Client) LegacyGetJRD(resource *Resource) (*jrd.JRD, error) {

	template, err := self.LegacyGetResourceJRDTemplateURL(resource.Domain)
	if err != nil {
		return nil, err
	}

	log.Printf("template: %s", template)

	jrdURL := strings.Replace(template, "{uri}", url.QueryEscape(resource.AsURIString()), 1)

	log.Printf("User JRD URL: %s", jrdURL)

	resourceJRD, err := self.findJRD([]string{jrdURL})
	if err != nil {
		return nil, err
	}

	return resourceJRD, nil
}
