package astrometry

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"telescope/util"
	"testing"
)

// Tests

func TestErrorResponseStruct(t *testing.T) {
	var structure = ErrorResponse{}
	err := json.Unmarshal(util.GetTestData(errorFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestLoginResponseStruct(t *testing.T) {
	var structure = LoginResponse{}
	err := json.Unmarshal(util.GetTestData(loginFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestUploadResponseStruct(t *testing.T) {
	var structure = UploadResponse{}
	err := json.Unmarshal(util.GetTestData(uploadFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestSubmissionStatusStruct(t *testing.T) {
	var structure = SubmissionStatus{}
	err := json.Unmarshal(util.GetTestData(subStatusFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestJobStatusStruct(t *testing.T) {
	var structure = JobStatus{}
	err := json.Unmarshal(util.GetTestData(jobStatusFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestCalibrationStruct(t *testing.T) {
	var structure = Calibration{}
	err := json.Unmarshal(util.GetTestData(calibrationFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestTaggedObjectsStruct(t *testing.T) {
	var structure = TaggedObjects{}
	err := json.Unmarshal(util.GetTestData(taggedObjectsFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestKnownObjectsStruct(t *testing.T) {
	var structure = KnownObjects{}
	err := json.Unmarshal(util.GetTestData(knownObjectsFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestAnnotationsListStruct(t *testing.T) {
	var structure = AnnotationsList{}
	err := json.Unmarshal(util.GetTestData(annotationsFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}

func TestJobResultsStruct(t *testing.T) {
	var structure = JobResults{}
	err := json.Unmarshal(util.GetTestData(jobResultsFile), &structure)
	if err != nil {
		t.Fatalf("testdata/struct failure: %s", err)
	}
	assert.NotNil(t, structure)
}
