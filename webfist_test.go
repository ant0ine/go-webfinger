package webfinger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/ant0ine/go-webfinger/jrd"
)

var (
	wfMux      *http.ServeMux
	wfServer   *httptest.Server
	wfTestHost string
)

func webFistSetup() {
	setup()

	wfMux = http.NewServeMux()
	wfServer = httptest.NewTLSServer(wfMux)
	u, _ := url.Parse(wfServer.URL)
	wfTestHost = u.Host

	client.WebFistServer = wfTestHost
}

func webFistTearDown() {
	teardown()
	wfServer.Close()
}

func TestWebFistLookup(t *testing.T) {
	webFistSetup()
	defer webFistTearDown()

	mux.HandleFunc("/webfinger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/jrd+json")
		fmt.Fprint(w, `{"subject":"bob@example.com"}`)
	})

	// simulate WebFist protocol
	wfMux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		resource := r.FormValue("resource")
		if want := "acct:bob@" + testHost; resource != want {
			t.Errorf("Requested resource: %v, want %v", resource, want)
		}
		w.Header().Add("content-type", "application/jrd+json")
		fmt.Fprint(w, `{
			"links": [{
				"rel": "http://webfist.org/spec/rel",
				"href": "`+server.URL+`/webfinger.json"
			}]
		}`)
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

func TestWebFistLookup_noLink(t *testing.T) {
	webFistSetup()
	defer webFistTearDown()

	wfMux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/jrd+json")
		fmt.Fprint(w, `{}`)
	})

	_, err := client.Lookup("bob@"+testHost, nil)
	if err == nil {
		t.Errorf("Expected webfist error.")
	}
}

func TestWebFistLookup_invalidLink(t *testing.T) {
	webFistSetup()
	defer webFistTearDown()

	wfMux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/jrd+json")
		fmt.Fprint(w, `{
			"links": [{
				"rel": "http://webfist.org/spec/rel",
				"href": "%"
			}]
		}`)
	})

	_, err := client.Lookup("bob@"+testHost, nil)
	if err == nil {
		t.Errorf("Expected webfist error.")
	}
}
