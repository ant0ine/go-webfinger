// Package webfinger provides a simple client implementation of the WebFinger
// protocol.
//
// It is a work in progress, the API is not frozen.
// We're trying to catchup with the last draft of the protocol:
// http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-14
// and to support the http://webfist.org
//
// Example:
//
//  package main
//
//  import (
//          "fmt"
//          "github.com/ant0ine/go-webfinger"
//          "os"
//  )
//
//  func main() {
//          email := os.Args[1]
//
//          client := webfinger.NewClient(nil)
//
//          resource, err := webfinger.MakeResource(email)
//          if err != nil {
//                  panic(err)
//          }
//
//          jrd, err := client.GetJRD(resource)
//          if err != nil {
//                  fmt.Println(err)
//                  return
//          }
//
//          fmt.Printf("JRD: %+v", jrd)
//  }
package webfinger

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-webfinger/jrd"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Resource represents a WebFinger resource.
type Resource struct {
	Local  string
	Domain string
}

// MakeResource constructs a WebFinger resource for the provided email string.
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

// AsURIString returns the resource as an URI string (eg: acct:user@domain).
func (self *Resource) AsURIString() string {
	return fmt.Sprintf("acct:%s@%s", self.Local, self.Domain)
}

// JRDURL returns the WebFinger URL that points to the JRD data for this resource.
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

// A Client is a WebFinger client.
type Client struct {
	// HTTP client used to perform WebFinger lookups.
	client *http.Client
}

// NewClient returns a new WebFinger client.  If a nil http.Client is provied,
// http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{client: httpClient}
}

// GetJRDPart returns the JRD for the specified resource, with the ability to
// specify which "rel" links to include.
func (self *Client) GetJRDPart(resource *Resource, rels []string) (*jrd.JRD, error) {

	log.Printf("Trying to get WebFinger JRD data for: %s", resource.AsURIString())

	resourceJRD, err := self.fetchJRD(resource.JRDURL(rels))
	if err != nil {
		return nil, err
	}

	return resourceJRD, nil
}

// GetJRD returns the JRD data for this resource.
// It follows redirect, and retries with http if https is not available.
func (self *Client) GetJRD(resource *Resource) (*jrd.JRD, error) {
	return self.GetJRDPart(resource, nil)
}

func (self *Client) fetchJRD(jrdURL *url.URL) (*jrd.JRD, error) {
	// TODO verify signature if not https
	// TODO extract http cache info

	// Get follows up to 10 redirects
	log.Printf("GET %s", jrdURL.String())
	res, err := self.client.Get(jrdURL.String())
	if err != nil {
		// retry with http instead of https
		if strings.Contains(err.Error(), "connection refused") {
			jrdURL.Scheme = "http"
			log.Printf("GET %s", jrdURL.String())
			res, err = self.client.Get(jrdURL.String())
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
	if strings.Contains(ct, "application/jrd+json") ||
		strings.Contains(ct, "application/json") {
		parsed, err := jrd.ParseJRD(content)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}

	return nil, errors.New(fmt.Sprintf("invalid content-type: %s", ct))
}
