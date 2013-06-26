package xrd

import (
	"testing"
)

func TestParseXRD(t *testing.T) {

	// from spec http://tools.ietf.org/html/rfc6415#appendix-A
	blob := `
<?xml version='1.0' encoding='UTF-8'?>
<XRD xmlns='http://docs.oasis-open.org/ns/xri/xrd-1.0'
 xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'>

<Subject>http://blog.example.com/article/id/314</Subject>
<Expires>2010-01-30T09:30:00Z</Expires>

<Alias>http://blog.example.com/cool_new_thing</Alias>
<Alias>http://blog.example.com/steve/article/7</Alias>

<Property type='http://blgx.example.net/ns/version'>1.2</Property>
<Property type='http://blgx.example.net/ns/version'>1.3</Property>
<Property type='http://blgx.example.net/ns/ext' xsi:nil='true' />

<Link rel='author' type='text/html'
    href='http://blog.example.com/author/steve'>
<Title>About the Author</Title>
<Title xml:lang='en-us'>Author Information</Title>
<Property type='http://example.com/role'>editor</Property>
</Link>

<Link rel='author' href='http://example.com/author/john'>
<Title>The other guy</Title>
<Title>The other author</Title>
</Link>
<Link rel='copyright'
    template='http://example.com/copyright?id={uri}' />
</XRD>`
	obj, err := ParseXRD([]byte(blob))
	if err != nil {
		t.Fatal(err)
	}
	if obj.Subject != "http://blog.example.com/article/id/314" {
		t.Error()
	}
	if obj.Property[0].Type != "http://blgx.example.net/ns/version" {
		t.Error()
	}
	if obj.Property[0].Value != "1.2" {
		t.Error()
	}
	if obj.GetLinkByRel("copyright") == nil {
		t.Error()
	}
	if obj.GetLinkByRel("copyright").Template != "http://example.com/copyright?id={uri}" {
		t.Error()
	}
	if obj.GetLinkByRel("author").Title[1].Value != "Author Information" {
		t.Error()
	}
	if obj.GetLinkByRel("author").Property[0].Value != "editor" {
		t.Error()
	}

	jrdObj := obj.ConvertToJRD()

	if jrdObj.Subject != "http://blog.example.com/article/id/314" {
		t.Error()
	}
	if jrdObj.Properties["http://blgx.example.net/ns/version"] != "1.3" {
		t.Error()
	}
	if jrdObj.GetLinkByRel("copyright") == nil {
		t.Error()
	}
	if jrdObj.GetLinkByRel("copyright").Template != "http://example.com/copyright?id={uri}" {
		t.Error()
	}
	if jrdObj.GetLinkByRel("author").Titles["default"] != "About the Author" {
		t.Error()
	}
	if jrdObj.GetLinkByRel("author").Properties["http://example.com/role"] != "editor" {
		t.Error()
	}
}
