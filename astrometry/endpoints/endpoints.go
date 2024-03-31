package endpoints

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func (e Endpoint) Body(key string) []byte {
	switch e {
	case 0:
		// Login
		return []byte(fmt.Sprintf(`request-json={"apikey": "%s"}`, key))

	case 1:
		// Upload URL
		return []byte(
			fmt.Sprintf(
				`request-json={"publicly_visible": "n", "allow_modifications": "n", "session": "%s", "allow_commercial_use": "n"}`,
				key,
			),
		)

	case 2:
		// Upload File
		return []byte(
			fmt.Sprintf(
				`{"publicly_visible": "n", "allow_modifications": "n", "session": "%s", "allow_commercial_use": "n"}`,
				key,
			),
		)

	default:
		return nil
	}
}

func (e Endpoint) Request(baseURL, keyOrID, filePath string) (*http.Request, error) {
	switch e {
	case 0, 1:
		// Login, Upload URL
		req, err := http.NewRequest(e.Method(), baseURL+e.URL(), bytes.NewBuffer(e.Body(keyOrID)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil

	case 2:
		// Upload File
		var reqBody bytes.Buffer
		mw := multipart.NewWriter(&reqBody)
		// Form field
		formField, err := mw.CreateFormField("request-json")
		if err != nil {
			return nil, err
		}
		if _, err = formField.Write(e.Body(keyOrID)); err != nil {
			return nil, err
		}
		// Open file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)
		// File field
		fileField, err := mw.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fileField, file); err != nil {
			return nil, err
		}
		err = mw.Close()
		if err != nil {
			return nil, err
		}
		// Create request
		req, err := http.NewRequest(e.Method(), baseURL+e.URL(), &reqBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())
		// Return
		return req, nil

	case 3, 4, 5, 6, 7, 8, 9:
		return http.NewRequest(e.Method(), baseURL+fmt.Sprintf(e.URL(), keyOrID), nil)
	default:
		return http.NewRequest(e.Method(), baseURL+e.URL(), nil)
	}
}
