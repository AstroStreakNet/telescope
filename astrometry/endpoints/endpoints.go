package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
)

// Type

type Endpoint int

// Enum

const (
	Login Endpoint = iota
	UploadUrl
	UploadFile
	SubmissionStatus
	JobStatus
	Calibration
	TaggedObjects
	KnownObjects
	Annotations
	JobResults
)

// Enum methods

func (e Endpoint) Method() string {
	switch e {
	case 0, 1, 2:
		return "POST"
	default:
		return "GET"
	}
}

func (e Endpoint) URL() string {
	switch e {
	case 0:
		return "/login"
	case 1:
		return "/url_upload"
	case 2:
		return "/upload"
	case 3:
		return "/submissions/%s"
	case 4:
		return "/jobs/%s"
	case 5:
		return "/jobs/%s/Calibration"
	case 6:
		return "/jobs/%s/machine_tags"
	case 7:
		return "/jobs/%s/objects_in_field"
	case 8:
		return "/jobs/%s/annotations"
	case 9:
		return "/jobs/%s/info"
	default:
		return ""
	}
}

func (e Endpoint) Body(key string) *bytes.Buffer {
	switch e {
	case 0:
		// Define arguments for request
		return bytes.NewBuffer([]byte(fmt.Sprintf(`request-json={"apikey": "%s"}`, key)))
	case 1, 2:
		return bytes.NewBuffer(
			[]byte(
				fmt.Sprintf(
					`request-json={"session": "%s", "allow_commercial_use": "n", "allow_modifications": "n", "publicly_visible": "n"}`,
					key,
				),
			),
		)
	default:
		return nil
	}
}

func (e Endpoint) Request(baseURL, keyOrID string) (*http.Request, error) {
	switch e {
	case 0, 1, 2:
		req, err := http.NewRequest(e.Method(), baseURL+e.URL(), e.Body(keyOrID))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	case 3, 4, 5, 6, 7, 8, 9:
		return http.NewRequest(e.Method(), baseURL+fmt.Sprintf(e.URL(), keyOrID), nil)
	default:
		return http.NewRequest(e.Method(), baseURL+e.URL(), nil)
	}
}
