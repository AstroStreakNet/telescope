package astrometry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"telescope/util"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func testSetupEndpoints() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(server.URL, "fakeKeyForAPI")

	return server.Close
}

func testSetupEndpointsError(endpoint string) func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(server.URL, "fakeKeyForAPI")

	// Define endpoint handler on test server
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(errorFile))
		if wErr != nil {
			panic(fmt.Sprintf("test server response write failure: %s", wErr))
		}
	})

	return server.Close
}

func TestSignIn(t *testing.T) {

}
