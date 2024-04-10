package astrometry

import (
	"encoding/json"
	"fmt"
	"github.com/AstroStreakNet/telescope/util"
	"github.com/stretchr/testify/assert"
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

func testSetupEndpointsError(endpoint string) func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(server.URL, "fakeKeyForAPI")
	client.SessionKey = "fakeKeyForSession"

	// Define endpoints handler on test server
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString("./testdata/error.json"))
		if wErr != nil {
			panic(fmt.Sprintf("test server response write failure: %s", wErr))
		}
	})

	return server.Close
}

func TestConnect(t *testing.T) {
	endTest := testSetup()
	defer endTest()

	mux.HandleFunc(Login.URL(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, util.GetTestDataString("./testdata/login.json")); err != nil {
			panic(fmt.Sprintf("test server response write failure: %s", err))
		}
	})

	req, err := Login.GetRequest(LoginOptions(client.baseURL, client.apiKey))
	if err != nil {
		t.Fatal(err)
	}
	resp := LoginResponse{}
	err = client.sendRequest(req, &resp)
	if err != nil {
		t.Fatal(err)
	}

	ref := LoginResponse{}
	err = json.Unmarshal(util.GetTestData("./testdata/login.json"), &ref)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", ref)
	t.Logf("%v", resp)
	assert.EqualValues(t, ref, resp)
}
