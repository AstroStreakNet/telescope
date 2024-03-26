package astrometry

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func testSetup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(server.URL, "fakeKeyForAPI")

	return server.Close
}

func TestSignIn(t *testing.T) {

}
