package webfinger

import (
	"errors"
	"github.com/ant0ine/go-webfinger/jrd"
	"log"
	"net/url"
	"strings"
)

// Build a serie well known host JRD URLs from the domain
func HostJRDURLs(domain string) []string {
	return []string{
		// first JRD implementation
		"https://" + domain + "/.well-known/host-meta.json",
		// orignal spec: https://code.google.com/p/webfinger/wiki/WebFingerProtocol
		"https://" + domain + "/.well-known/host-meta",
	}
}

// Try to call FetchJRD on each url until a successful response.
func FindJRD(urls []string) (*jrd.JRD, error) {
	for _, try := range urls {
		try_obj, err := url.Parse(try)
		if err != nil {
			log.Print(err)
			continue
		}
		obj, err := FetchJRD(try_obj)
		if err != nil {
			log.Print(err)
			continue
		}
		return obj, nil
	}
	return nil, errors.New("JRD not found")
}

// Given a domain, this method gets the host meta JRD data,
// and returns the LRDD user JRD template URL.
func GetUserJRDTemplateURL(domain string) (string, error) {
	// TODO implement heavy HTTP cache around this

	urls := HostJRDURLs(domain)

	host_jrd, err := FindJRD(urls)
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

func (self *Resource) GetJRDCompat() (*jrd.JRD, error) {

	template, err := GetUserJRDTemplateURL(self.Domain)
	if err != nil {
		return nil, err
	}

	log.Printf("template: %s", template)

	jrd_url := strings.Replace(template, "{uri}", url.QueryEscape(self.AsURIString()), 1)

	log.Printf("User JRD URL: %s", jrd_url)

	resource_jrd, err := FindJRD([]string{jrd_url})
	if err != nil {
		return nil, err
	}

	return resource_jrd, nil
}
