package astrometry

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"telescope/util"
	"testing"
)

func TestLogin(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	const testFile = "/astrometry/endpoints/login.json"

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// TODO check request structure as part of test

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, util.GetTestDataString(testFile))
	})

	response, loginErr := client.login()
	if loginErr != nil {
		t.Fatal(loginErr)
	}

	var reference = LoginResponse{}
	decodeErr := json.Unmarshal(util.GetTestData(testFile), &reference)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, response.Status, reference.Status, "Response should match json file data correctly")
	assert.Equal(t, response.Message, reference.Message, "Response should match json file data correctly")
	assert.Equal(t, response.Session, reference.Session, "Response should match json file data correctly")
}

func TestLoginErrorResponse(t *testing.T) {
	testEnd := testSetup()
	defer testEnd()

	const testFile = "/astrometry/endpoints/error.json"

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, util.GetTestDataString(testFile))
	})

	_, loginErr := client.login()
	assert.Error(t, loginErr, "error should occur as response was an error message")
}

func TestUpload(t *testing.T) {
}
