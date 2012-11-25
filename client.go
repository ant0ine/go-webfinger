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
//      resource, err := webfinger.MakeResource("user@host")
//	if err != nil {
//		panic(err)
//	}
//	jrd, err := resource.GetJRD()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("JRD: %+v", ujrd)
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

// WebFinger Resource
type Resource struct {
	Local  string
	Domain string
}

// Parse the email string and return a *Resource
func MakeResource(email string) (*Resource, error) {
	// TODO validate address, see http://www.ietf.org/rfc/rfc2822.txt
	// TODO accept an email address URI
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 {
		return nil, errors.New("not a valid email")
	}
	return &Resource{
		Local:  parts[0],
		Domain: parts[1],
	}, nil
}

// Return the resource as an URI (ie: acct:user@domain)
func (self *Resource) AsURI() string {
	return fmt.Sprintf("acct:%s@%s", self.Local, self.Domain)
}

// Generate the WebFinger URLs that can point to the resource JRD
func (self *Resource) JRDURLs() []string {
	// TODO support the rel query string parameter
	uri := url.QueryEscape(self.AsURI())
	return []string{
		"https://" + self.Domain + "/.well-known/webfinger?resource=" + uri,
		"http://" + self.Domain + "/.well-known/webfinger?resource=" + uri,
	}
}

// Try to get the JRD data for this resource
func (self *Resource) GetJRD() (*jrd.JRD, error) {
	// TODO support the rel query string parameter

	log.Printf("Trying to get WebFinger JRD data for: %s", self.AsURI())

	resource_jrd, err := FindJRD(self.JRDURLs())
	if err != nil {
		// Fallback to the original WebFinger spec
		log.Print(err)
		resource_jrd, err = self.GetJRDCompat()
		if err != nil {
			return nil, err
		}
	}

	if resource_jrd.Subject != self.AsURI() {
		return nil, errors.New(fmt.Sprintf("JRD Subject does not match the resource: %s", self.AsURI()))
	}

	return resource_jrd, nil
}

// Given an URL, get and parse the JRD.
// [Compat Note] If the payload is in XRD, this method parses it
// and converts it to JRD.
func FetchJRD(url string) (*jrd.JRD, error) {
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

// Try to call FetchJRD on each url until a successful response.
func FindJRD(urls []string) (*jrd.JRD, error) {
	for _, url := range urls {
		log.Printf("Fetching Host JRD URL: %s", url)
		obj, err := FetchJRD(url)
		if err == nil {
			return obj, nil
		}
		log.Print(err)
	}
	return nil, errors.New("JRD not found")
}
