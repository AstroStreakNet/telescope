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

	subID = "12345"
	jobID = "54321"

	sessionKey = "api"

	fileToUpload = ""
)

func TestLogin(t *testing.T) {
	testEnd := testSetupEndpoints()
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
	testEnd := testSetupEndpointsError(loginEP)
	defer testEnd()

	_, respErr := client.login()
	assert.Error(t, respErr, assertErrorMsg)
}

func TestUpload(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(uploadEP, func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(uploadFile))
		if wErr != nil {
			t.Fatal("test server response write failure")
		}
	})

	// Send request to test server
	response, respErr := client.upload(fileToUpload)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = UploadResponse{}
	decodeErr := json.Unmarshal(util.GetTestData(uploadFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	// Check if response and reference are equal
	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestUploadError(t *testing.T) {
	testEnd := testSetupEndpointsError(uploadEP)
	defer testEnd()

	_, respErr := client.upload(fileToUpload)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetSubmissionStatus(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(subStatusEP, subID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(subStatusFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getSubmissionStatus(subID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = SubmissionStatus{}
	decodeErr := json.Unmarshal(util.GetTestData(subStatusFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetSubmissionStatusError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(subStatusEP, subID))
	defer testEnd()

	_, respErr := client.getSubmissionStatus(subID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetJobStatus(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(jobStatusEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(jobStatusFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getJobStatus(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = JobStatus{}
	decodeErr := json.Unmarshal(util.GetTestData(jobStatusFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetJobStatusError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(jobStatusEP, jobID))
	defer testEnd()

	_, respErr := client.getSubmissionStatus(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetCalibration(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(calibrationEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(calibrationFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getCalibration(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = Calibration{}
	decodeErr := json.Unmarshal(util.GetTestData(calibrationFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetCalibrationError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(calibrationEP, jobID))
	defer testEnd()

	_, respErr := client.getCalibration(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetTaggedObjects(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(taggedObjectsEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(taggedObjectsFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getTaggedObjects(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = TaggedObjects{}
	decodeErr := json.Unmarshal(util.GetTestData(taggedObjectsFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetTaggedObjectsError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(taggedObjectsEP, jobID))
	defer testEnd()

	_, respErr := client.getTaggedObjects(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetKnownObjects(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(knownObjectsEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(knownObjectsFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getKnownObjects(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = KnownObjects{}
	decodeErr := json.Unmarshal(util.GetTestData(knownObjectsFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetKnownObjectsError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(knownObjectsEP, jobID))
	defer testEnd()

	_, respErr := client.getKnownObjects(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetAnnotations(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(annotationsEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(annotationsFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getAnnotations(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = AnnotationsList{}
	decodeErr := json.Unmarshal(util.GetTestData(annotationsFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetAnnotationsError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(annotationsEP, jobID))
	defer testEnd()

	_, respErr := client.getAnnotations(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}

func TestGetJobResults(t *testing.T) {
	testEnd := testSetupEndpoints()
	defer testEnd()

	// Define endpoint handler on test server
	mux.HandleFunc(fmt.Sprintf(jobResultsEP, jobID), func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, wErr := fmt.Fprint(w, util.GetTestDataString(jobResultsFile))
		if wErr != nil {
			t.Fatalf("test server response write failure: %s", wErr)
		}
	})

	// Send request to test server
	response, respErr := client.getJobResults(jobID)
	assert.Nil(t, respErr, assertNilRespMsg)

	// Retrieve json to create a reference struct for comparison
	var reference = JobResults{}
	decodeErr := json.Unmarshal(util.GetTestData(jobResultsFile), &reference)
	if decodeErr != nil {
		t.Fatalf("testdata/struct failure: %s", decodeErr)
	}

	assert.EqualValues(t, response, &reference, assertEqualMsg)
}

func TestGetJobResultsError(t *testing.T) {
	testEnd := testSetupEndpointsError(fmt.Sprintf(jobResultsEP, jobID))
	defer testEnd()

	_, respErr := client.getJobResults(jobID)
	assert.Error(t, respErr, assertErrorMsg)
}
