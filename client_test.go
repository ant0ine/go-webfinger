package webfinger

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/ant0ine/go-webfinger/jrd"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server

	// testHost is the hostname and port of the local running test server.
	testHost string

	// client is the WebFinger client being tested.
	client *Client
)

// setup a local HTTP server for testing
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewTLSServer(mux)
	u, _ := url.Parse(server.URL)
	testHost = u.Host

	// for testing, use an HTTP client which doesn't check certs
	client = NewClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
}

func teardown() {
	server.Close()
}

func TestResource_Parse(t *testing.T) {
	// URL with host
	r, err := Parse("http://example.com/")
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
	want := &Resource{Scheme: "http", Host: "example.com", Path: "/"}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("Parsed resource: %#v, want %#v", r, want)
	}

	// email-like identifier
	r, err = Parse("bob@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
	want = &Resource{Scheme: "acct", Opaque: "bob@example.com"}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("Parsed resource: %#v, want %#v", r, want)
	}
}

func TestResource_Parse_error(t *testing.T) {
	_, err := Parse("example.com")
	if err == nil {
		t.Error("Expected parse error")
	}

	_, err = Parse("%")
	if err == nil {
		t.Error("Expected parse error")
	}
}

func TestResource_WebFingerHost(t *testing.T) {
	// URL with host
	r, _ := Parse("http://example.com/")
	if got, want := r.WebFingerHost(), "example.com"; got != want {
		t.Errorf("WebFingerHost() returned: %#v, want %#v", got, want)
	}

	// email-like identifier
	r, _ = Parse("bob@example.com")
	if got, want := r.WebFingerHost(), "example.com"; got != want {
		t.Errorf("WebFingerHost() returned: %#v, want %#v", got, want)
	}

	// mailto URL
	r, _ = Parse("mailto:bob@example.com")
	if got, want := r.WebFingerHost(), "example.com"; got != want {
		t.Errorf("WebFingerHost() returned: %#v, want %#v", got, want)
	}

	// URL with no host
	r, _ = Parse("file:///example")
	if got, want := r.WebFingerHost(), ""; got != want {
		t.Errorf("WebFingerHost() returned: %#v, want %#v", got, want)
	}
}

func TestResource_JRDURL(t *testing.T) {
	r, _ := Parse("bob@example.com")
	got := r.JRDURL("", nil)
	want, _ := url.Parse("https://example.com/.well-known/webfinger?" +
		"resource=acct%3Abob%40example.com")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JRDURL() returned: %#v, want %#v", got, want)
	}

	r, _ = Parse("http://example.com/")
	got = r.JRDURL("example.net", []string{"a", "b"})
	// sadly, we have to compare each URL component individually because the
	// order of query string values is unpredictable
	if want := "https"; got.Scheme != want {
		t.Errorf("JRDURL() returned scheme: %#v, want %#v", got.Scheme, want)
	}
	if want := "example.net"; got.Host != want {
		t.Errorf("JRDURL() returned host: %#v, want %#v", got.Host, want)
	}
	if want := "/.well-known/webfinger"; got.Path != want {
		t.Errorf("JRDURL() returned path: %#v, want %#v", got.Path, want)
	}
	if want := []string{"http://example.com/"}; reflect.DeepEqual(got.Query().Get("resource"), want) {
		t.Errorf("JRDURL() returned query resource: %#v, want %#v", got.Query().Get("resource"), want)
	}
	if want := []string{"a", "b"}; reflect.DeepEqual(got.Query().Get("rel"), want) {
		t.Errorf("JRDURL() returned query rel: %#v, want %#v", got.Query().Get("rel"), want)
	}
}

func TestResource_String(t *testing.T) {
	r, _ := Parse("bob@example.com")
	if got, want := r.String(), "acct:bob@example.com"; got != want {
		t.Errorf("String() returned: %#v, want %#v", got, want)
	}
}

func TestLookup(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		resource := r.FormValue("resource")
		if want := "acct:bob@" + testHost; resource != want {
			t.Errorf("Requested resource: %v, want %v", resource, want)
		}
		w.Header().Add("content-type", "application/jrd+json")
		fmt.Fprint(w, `{"subject":"bob@example.com"}`)
	})

	JRD, err := client.Lookup("bob@"+testHost, nil)
	if err != nil {
		t.Errorf("Unexpected error lookup up webfinger: %#v", err)
	}
	want := &jrd.JRD{Subject: "bob@example.com"}
	if !reflect.DeepEqual(JRD, want) {
		t.Errorf("Lookup returned %#v, want %#v", JRD, want)
	}
}

func TestLookup_parseError(t *testing.T) {
	// use default client here, just to make sure that gets tested
	_, err := Lookup("bob", nil)
	if err == nil {
		t.Error("Expected parse error")
	}
}

func TestLookup_404(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Lookup("bob@"+testHost, nil)
	if err == nil {
		t.Error("Expected error")
	}
}
