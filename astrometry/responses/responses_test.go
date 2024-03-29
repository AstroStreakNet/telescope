package responses

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"telescope/util"
	"testing"
)

func TestResponseDecoding(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		response interface{}
	}{
		{"error response", "../testdata/error.json", Error{}},
		{"login response", "../testdata/login.json", Login{}},
		{"upload response", "../testdata/upload.json", Upload{}},
		{"submission status", "../testdata/submission_status.json", SubmissionStatus{}},
		{"job status", "../testdata/job_status.json", JobStatus{}},
		{"calibration", "../testdata/calibration.json", Calibration{}},
		{"tagged objects", "../testdata/tagged_objects.json", TaggedObjects{}},
		{"known objects", "../testdata/known_objects.json", KnownObjects{}},
		{"annotations", "../testdata/annotations.json", Annotations{}},
		{"job results", "../testdata/job_results.json", JobResults{}},
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
