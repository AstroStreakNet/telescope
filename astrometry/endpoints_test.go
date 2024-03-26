package astrometry

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	loginFile   = "/astrometry/login.json"
	errorFile   = "/astrometry/error.json"
	subStatFile = "/astrometry/submission_status.json"
	jobStatFile = "/astrometry/job_status.json"

	subID = "12345"
	jobID = "12345"
)

func TestLogin(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(loginEP, func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(loginFile))
		if wErr != nil {
			t.Fatal("test server response write failure")
		}
	})

	// Send request to test server
	response, respErr := client.login()
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = LoginResponse{}
	decodeErr := json.Unmarshal(util.GetTestData(loginFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	// Check if response and reference are equal
	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestLoginError(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	mux.HandleFunc(loginEP, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(errorFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	_, respErr := client.login()
	assert.Error(t, respErr, assertErrorMsg)
}

func TestUpload(t *testing.T) {
}

func TestUploadError(t *testing.T) {
}

func TestGetSubmissionStatus(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(subStatusEP, subID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(subStatFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getSubmissionStatus(subID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = SubmissionStatus{}
	decodeErr := json.Unmarshal(util.GetTestData(subStatFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetSubmissionStatusError(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(subStatusEP, subID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(errorFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	_, respErr := client.getSubmissionStatus(subID)
	assert.Error(t, respErr, "error should occur as response was an error message")
}

func TestGetJobStatus(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()
	// constants for test

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(jobStatusEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(jobStatFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getJobStatus(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = JobStatus{}
	decodeErr := json.Unmarshal(util.GetTestData(jobStatFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}
