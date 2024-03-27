package astrometry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"telescope/util"
	"testing"
)

// Assertion messages
const (
	assertEqualMsg   = "Response should match json file data correctly"
	assertErrorMsg   = "error should occur as response was an error message"
	assertNilRespMsg = "Shouldn't be any errors in request creation and response decoding"
)

// Test data file paths and ids
const (
	errorFile         = "/astrometry/error.json"
	loginFile         = "/astrometry/login.json"
	uploadFile        = "/astrometry/upload.json"
	subStatusFile     = "/astrometry/submission_status.json"
	jobStatusFile     = "/astrometry/job_status.json"
	calibrationFile   = "/astrometry/calibration.json"
	taggedObjectsFile = "/astrometry/tagged_objects.json"
	knownObjectsFile  = "/astrometry/known_objects.json"
	annotationsFile   = "/astrometry/annotations.json"
	jobResultsFile    = "/astrometry/job_results.json"
	fileToUpload      = "/astrometry/file_to_upload.txt"

	subID = "12345"
	jobID = "54321"
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
	client.SessionKey = "fakeKeyForSession"

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
