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
//	fmt.Printf("JRD: %+v", jrd)
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
	// TODO support mailto: http:  <= rework that
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 {
		return nil, errors.New("not a valid email")
	}
	return &Resource{
		Local:  parts[0],
		Domain: parts[1],
	}, nil
}

// Return the resource as an URI string (eg: acct:user@domain)
func (self *Resource) AsURIString() string {
	return fmt.Sprintf("acct:%s@%s", self.Local, self.Domain)
}

// Generate the WebFinger URL that points to the JRD data for this resource.
func (self *Resource) JRDURL(rels []string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   self.Domain,
		Path:   "/.well-known/webfinger",
		RawQuery: url.Values{
			"resource": []string{self.AsURIString()},
			"rel":      rels,
		}.Encode(),
	}
}

// Get the JRD data for this resource with the ability to specify which "rel" links to include.
func (self *Resource) GetJRDPart(rels []string) (*jrd.JRD, error) {

	log.Printf("Trying to get WebFinger JRD data for: %s", self.AsURIString())

	resource_jrd, err := FetchJRD(self.JRDURL(rels))
	if err != nil {
		log.Print(err)
		log.Print("Fallback to the original WebFinger spec")
		resource_jrd, err = self.GetJRDCompat()
		if err != nil {
			return nil, err
		}
	}

	// verify the subject
	if resource_jrd.Subject != self.AsURIString() {
		return nil, errors.New(
			fmt.Sprintf(
				"JRD Subject does not match the resource: %s",
				self.AsURIString(),
			),
		)
	}

	return resource_jrd, nil
}

// Get the JRD data for this resource.
func (self *Resource) GetJRD() (*jrd.JRD, error) {
	return self.GetJRDPart(nil)
}

// Given an URL, get and parse the JRD.
// It follows redirect, and retries with http if https is not available.
// [Compat Note] If the payload is in XRD, this method parses it
// and converts it to JRD.
func FetchJRD(jrd_url *url.URL) (*jrd.JRD, error) {
	// TODO verify signature if not https
	// TODO extract http cache info

	// Get follows up to 10 redirects
	log.Printf("GET %s", jrd_url.String())
	res, err := http.Get(jrd_url.String())
	if err != nil {
		// retry with http instead of https
		if strings.Contains(err.Error(), "connection refused") {
			jrd_url.Scheme = "http"
			log.Printf("GET %s", jrd_url.String())
			res, err = http.Get(jrd_url.String())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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
