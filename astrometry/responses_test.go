package astrometry

import (
	"encoding/json"
	"github.com/AstroStreakNet/telescope/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseDecoding(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		response interface{}
	}{
		{"error response", "./testdata/error.json", ErrorResponse{}},
		{"login response", "./testdata/login.json", LoginResponse{}},
		{"upload response", "./testdata/upload.json", UploadResponse{}},
		{"submission status", "./testdata/submission_status.json", SubmissionStatusResponse{}},
		{"calibration", "./testdata/calibration.json", CalibrationResponse{}},
		{"annotations", "./testdata/annotations.json", AnnotationsResponse{}},
		{"job results", "./testdata/job_results.json", JobResultsResponse{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := json.Unmarshal(util.GetTestData(test.data), &test.response)
			if err != nil {
				t.Fatalf("%s failure: %s", test.name, err)
			}
			t.Logf("response values: %s", test.response)
			assert.NotNil(t, test.response)
		})
	}
}
