package astrometry

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
	Upload
	SubmissionStatus
	Calibration
	Annotations
	JobResults
)

// Enum methods

func (e Endpoint) URL() string {
	switch e {
	case Login:
		return "/login"
	case Upload:
		return "/upload"
	case SubmissionStatus:
		return "/submissions/%s"
	case Calibration:
		return "/jobs/%s/Calibration"
	case Annotations:
		return "/jobs/%s/annotations"
	case JobResults:
		return "/jobs/%s/info"
	default:
		return ""
	}
}

func (e Endpoint) Arguments() string {
	switch e {
	case Login:
		return `request-json={"apikey": "%s"`
	case Upload:
		return `{"publicly_visible": "n", "allow_modifications": "n", "session": "%s", "allow_commercial_use": "n"}`
	default:
		return ""
	}
}

func (e Endpoint) Method() string {
	switch e {
	case 0, 1, 2:
		return "POST"
	default:
		return "GET"
	}
}

func (e Endpoint) GetRequest(options *Options) (*http.Request, error) {
	switch e {

	case Login:
		return http.NewRequest(
			e.Method(),
			options.BaseURL+e.URL(),
			bytes.NewBuffer([]byte(fmt.Sprintf(e.Arguments(), options.ApiKey))),
		)

	case Upload:
		var reqBody bytes.Buffer
		// Instantiate multipart form creator
		mw := multipart.NewWriter(&reqBody)

		// Form field
		formField, err := mw.CreateFormField("request-json")
		if err != nil {
			return nil, err
		}
		if _, err = formField.Write(
			[]byte(fmt.Sprintf(e.Arguments(), options.SessionKey)),
		); err != nil {
			return nil, err
		}

		// Open file
		file, err := os.Open(options.FilePath)
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
		fileField, err := mw.CreateFormFile("file", filepath.Base(options.FilePath))
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
		req, err := http.NewRequest(e.Method(), options.BaseURL+e.URL(), &reqBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())

		// Return
		return req, nil

	case SubmissionStatus:
		return http.NewRequest(
			e.Method(),
			options.BaseURL+fmt.Sprintf(e.URL(), options.SubmissionID),
			nil,
		)

	case Calibration, JobResults:
		return http.NewRequest(
			e.Method(),
			options.BaseURL+fmt.Sprintf(e.URL(), options.JobID),
			nil,
		)

	default:
		return http.NewRequest(
			e.Method(),
			options.BaseURL+e.URL(),
			nil,
		)
	}
}

// Options struct

type Options struct {
	BaseURL      string
	ApiKey       string
	SessionKey   string
	SubmissionID int
	JobID        int
	FilePath     string
}

// Options factory functions

func LoginOptions(baseURL, apiKey string) *Options {
	return &Options{
		BaseURL: baseURL,
		ApiKey:  apiKey,
	}
}

func UploadOptions(baseURL, sessionKey, filePath string) *Options {
	return &Options{
		BaseURL:    baseURL,
		SessionKey: sessionKey,
		FilePath:   filePath,
	}
}

func SubmissionStatusOptions(baseURL string, submissionID int) *Options {
	return &Options{
		BaseURL:      baseURL,
		SubmissionID: submissionID,
	}
}

func CalibrationOptions(baseURL string, jobID int) *Options {
	return &Options{
		BaseURL: baseURL,
		JobID:   jobID,
	}
}

func JobResultsOptions(baseURL string, jobID int) *Options {
	return &Options{
		BaseURL: baseURL,
		JobID:   jobID,
	}
}
