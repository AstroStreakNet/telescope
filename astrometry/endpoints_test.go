package astrometry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndpointGetRequest(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		options  *Options
		method   string
		url      string
	}{
		{
			"login", Login,
			LoginOptions("test.test", "key"),
			"POST", "test.test/login",
		},
		{
			"upload", Upload,
			UploadOptions("test.test", "key", "./testdata/file_to_upload.txt"),
			"POST", "test.test/upload",
		},
		{
			"submission status", SubmissionStatus,
			SubmissionStatusOptions("test.test", 0),
			"GET", "test.test/submissions/0",
		},
		{
			"calibration", Calibration,
			CalibrationOptions("test.test", 0),
			"GET", "test.test/jobs/0/calibration",
		},
		{
			"annotations", Annotations,
			AnnotationsOptions("test.test", 0),
			"GET", "test.test/jobs/0/annotations",
		},
		{
			"job results", JobResults,
			JobResultsOptions("test.test", 0),
			"GET", "test.test/jobs/0/info",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := test.endpoint.GetRequest(test.options)
			assert.Nil(t, err)
			assert.Equal(t, test.method, req.Method)
			assert.Equal(t, test.url, req.URL.String())
		})
	}
}
