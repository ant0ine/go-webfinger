// Simple Client Implementation of WebFinger
//
// This is a work in progress, The WebFinger spec is changing every day anyway :)
//
// Some links to follow:
// https://code.google.com/p/webfinger/wiki/WebFingerProtocol
// https://groups.google.com/forum/#!forum/webfinger
// https://groups.google.com/d/topic/webfinger/zw-pCRGyuSo/discussion
// http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-04
package webfinger

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-webfinger/xrd"
	"log"
	"net/url"
	"strings"
	"io/ioutil"
	"net/http"
)

type EmailAddress struct { // XXX rename Id or WebFingerId ? it may not be an email address
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

// Build the well known host meta URL from the domain
func HostXRDURL(domain string) string { // XXX s/XRD/Meta/ ?
	// TODO return also the non-https URL
	// TODO return the JRD URL ?
	return "https://" + domain + "/.well-known/host-meta"
}

func GetXRD(url string) (*xrd.XRD, error) {
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

	parsed, err := xrd.ParseXRD(content)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// Given a domain, this method gets the host meta XRD data,
// and returns the LRDD user XRD template URL.
func GetUserXRDTemplateURL(domain string) (string, error) {
	// TODO implement heavy HTTP cache around this

	xrd_url := HostXRDURL(domain)

	log.Printf("Fetching Host XRD URL: %s", xrd_url)

	host_xrd, err := GetXRD(xrd_url)
	if err != nil {
		return "", err
	}

	template := host_xrd.LrddTemplate()
	if template == "" {
		return "", errors.New("cannot find the template in the XRD data")
	}

	return template, nil
}

// Try to discover the user XRD data from the email
func GetUserXRD(email string) (*xrd.XRD, error) {
        // TODO support the rel query string parameter

	address, err := MakeEmailAddress(email)
	if err != nil {
		return nil, err
	}

	log.Printf("Fetching WebFinger XRD info for: %s", address.AsURI())

	template, err := GetUserXRDTemplateURL(address.Domain)
	if err != nil {
		return nil, err
	}

	log.Printf("template: %s", template)

	xrd_url := strings.Replace(template, "{uri}", url.QueryEscape(address.AsURI()), 1)

	log.Printf("User XRD URL: %s", xrd_url)

	user_xrd, err := GetXRD(xrd_url)
	if err != nil {
		return nil, err
	}

	if user_xrd.Subject != address.AsURI() {
		return nil, errors.New("XRD Subject does not match the email")
	}

	return user_xrd, nil
}
