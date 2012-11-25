// Simple Client Implementation of WebFinger
//
// (This is a work in progress, the API is not frozen)
//
// This implementation tries to follow the last spec:
// http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-04
// And also tries to provide backwark compatibility with the original spec:
// https://code.google.com/p/webfinger/wiki/WebFingerProtocol
//
// Example:
//      import (
//          "fmt"
//          "github.com/ant0ine/go-webfinger"
//      )
//
// 	jrd, err := webfinger.GetUserJRD("user@host")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("User JRD: %+v", jrd)
package webfinger

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-webfinger/jrd"
	"github.com/ant0ine/go-webfinger/xrd"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type EmailAddress struct { // TODO rename Id or WebFingerId ? it may not be an email address
	Local  string
	Domain string
}

// Parse the email string and return an *EmailAddress
func MakeEmailAddress(email string) (*EmailAddress, error) {
	// TODO validate address, see http://www.ietf.org/rfc/rfc2822.txt
	// TODO accept an email address URI
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 {
		return nil, errors.New("not a valid email")
	}
	return &EmailAddress{
		Local:  parts[0],
		Domain: parts[1],
	}, nil
}

// Return the email address as an URI (ie: acct:user@domain)
func (self *EmailAddress) AsURI() string {
	return fmt.Sprintf("acct:%s@%s", self.Local, self.Domain)
}

// Build a serie well known host JRD URLs from the domain
// [Compat Note] This includes URLs from previous versions of the spec.
func HostJRDURLs(domain string) []string {
	return []string{
		// last spec: http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-04
		"https://" + domain + "/.well-known/webfinger",
		"http://" + domain + "/.well-known/webfinger",
		// first JRD implementation
		"https://" + domain + "/.well-known/host-meta.json",
		"http://" + domain + "/.well-known/host-meta.json",
		// orignal spec: https://code.google.com/p/webfinger/wiki/WebFingerProtocol
		"https://" + domain + "/.well-known/host-meta",
		"http://" + domain + "/.well-known/host-meta",
	}
}

// Given an URL, get and parse the JRD.
// [Compat Note] If the payload is in XRD, this method parses it
// and converts it to JRD.
func GetJRD(url string) (*jrd.JRD, error) {
	// TODO verify signature if not https
	// TODO extract http cache info

	// Get follows up to 10 redirects
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return nil, errors.New(res.Status)
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	ct := strings.ToLower(res.Header.Get("content-type"))
	if strings.Contains(ct, "application/json") {
		parsed, err := jrd.ParseJRD(content)
		if err != nil {
			return nil, err
		}
		return parsed, nil

	} else if strings.Contains(ct, "application/xrd+xml") ||
		strings.Contains(ct, "application/xml") ||
		strings.Contains(ct, "text/xml") {
		parsed, err := xrd.ParseXRD(content)
		if err != nil {
			return nil, err
		}
		return parsed.ConvertToJRD(), nil
	}

	return nil, errors.New(fmt.Sprintf("invalid content-type: %s", ct))
}

// Try to call GetJRD on each url until a successful response.
func FindJRD(urls []string) (*jrd.JRD, error) {
	for _, url := range urls {
		log.Printf("Fetching Host JRD URL: %s", url)
		obj, err := GetJRD(url)
		if err == nil {
			return obj, nil
		}
		log.Print(err)
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

	template := host_jrd.LrddTemplate()
	if template == "" {
		return "", errors.New("cannot find the template in the JRD data")
	}

	return template, nil
}

// Try to discover the user JRD data from the email
func GetUserJRD(email string) (*jrd.JRD, error) {
	// TODO support the rel query string parameter

	address, err := MakeEmailAddress(email)
	if err != nil {
		return nil, err
	}

	log.Printf("Fetching WebFinger JRD info for: %s", address.AsURI())

	template, err := GetUserJRDTemplateURL(address.Domain)
	if err != nil {
		return nil, err
	}

	log.Printf("template: %s", template)

	jrd_url := strings.Replace(template, "{uri}", url.QueryEscape(address.AsURI()), 1)

	log.Printf("User JRD URL: %s", jrd_url)

	user_jrd, err := GetJRD(jrd_url)
	if err != nil {
		return nil, err
	}

	if user_jrd.Subject != address.AsURI() {
		return nil, errors.New("JRD Subject does not match the email")
	}

	return user_jrd, nil
}
