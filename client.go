// Package webfinger provides a simple client implementation of the WebFinger
// protocol.
//
// (This is a work in progress, the API is not frozen)
//
// This implementation tries to follow the last spec:
// http://tools.ietf.org/html/draft-ietf-appsawg-webfinger-05
//
// And also tries to provide backwark compatibility with the original spec:
// https://code.google.com/p/webfinger/wiki/WebFingerProtocol
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

// Resource is a resource for which a WebFinger query can be issued.
type Resource url.URL

// Parse parses rawurl into a WebFinger Resource.  The rawurl should be an
// absolute URL, or an email-like identifier (e.g. "bob@example.com").
func Parse(rawurl string) (*Resource, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if parsed URL has no scheme but is email-like, treat it as an acct: URL.
	if u.Scheme == "" {
		parts := strings.SplitN(rawurl, "@", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("URL must be absolute, or an email address: %v", rawurl)
		}
		return Parse("acct:" + rawurl)
	}

	r := Resource(*u)
	return &r, nil
}

// WebFingerHost returns the default host for issuing WebFinger queries for
// this resource.  For Resource URLs with a host component, that value is used.
// For URLs that do not have a host component, the host is determined by other
// mains if possible (for example, the domain in the addr-spec of a mailto
// URL).  If the host cannot be determined from the URL, this value will be an
// empty string.
func (r *Resource) WebFingerHost() string {
	if r.Host != "" {
		return r.Host
	} else if r.Scheme == "acct" || r.Scheme == "mailto" {
		parts := strings.SplitN(r.Opaque, "@", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}

// String reassembles the Resource into a valid URL string.
func (r *Resource) String() string {
	u := url.URL(*r)
	return u.String()
}

// JRDURL returns the WebFinger query URL at the specified host for this
// resource.  If host is an empty string, the default host for the resource
// will be used, as returned from WebFingerHost().
func (r *Resource) JRDURL(host string, rels []string) *url.URL {
	if host == "" {
		host = r.WebFingerHost()
	}

	return &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/.well-known/webfinger",
		RawQuery: url.Values{
			"resource": []string{r.String()},
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

	log.Printf("Trying to get WebFinger JRD data for: %s", resource.String())

	resourceJRD, err := self.fetchJRD(resource.JRDURL("", rels))
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
