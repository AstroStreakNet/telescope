package astrometry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Client

type Client struct {
	apiKey     string
	baseURL    string
	SessionKey string
	httpClient *http.Client
}

// Client factory functions

func NewAstrometryClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "http://nova.astrometry.net/api", // Astrometry doesn't support https
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Private methods

// sendRequest is a Client method that sends off requests, checks if the response is an error message,
// and if not then decodes the data into a specified data structure
func (c *Client) sendRequest(req *http.Request, respStruct interface{}) error {
	// Send provided request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	// If no error in sending request than defer close
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("ReadCloser failure for response data: %s", err))
		}
	}(resp.Body)

	// Read bytes from response body
	var respBytes bytes.Buffer
	_, err = io.Copy(&respBytes, resp.Body)
	if err != nil {
		return err
	}
	// Create copy of bytes for error message detection.
	checkBytes := bytes.NewBuffer(respBytes.Bytes())

	// Error checking has to be done in this roundabout way due to Astrometry API not properly using http status codes,
	// it returns 200, http ok status code, regardless of whether the query was actually a success or not.
	// The response body has to be checked to determine if an error has occurred.
	// It is very annoying.
	var errorResponse = ErrorResponse{}
	err = json.NewDecoder(checkBytes).Decode(&errorResponse)
	if err != nil {
		return err
	} else if errorResponse.Status == "error" {
		return errors.New(fmt.Sprintf("error response from astrometry: %s", errorResponse.ErrorMessage))
	}

	// Decode response into desired data structure
	err = json.NewDecoder(&respBytes).Decode(&respStruct)
	if err != nil {
		return err
	}

	// Return nothing if no errors have occurred throughout whole process
	return nil
}

func (c *Client) sessionCheck(e error) error {
	// Check if session key has expired.
	if strings.Contains(e.Error(), "no session with key") {
		// If session key expired, login again
		_, err := c.Connect()
		if err != nil {
			// Session key was issue, failed to reconnect
			return err
		}
		// Session key was issue, reconnected
		return nil
	}
	// Session key wasn't issue
	return e
}

// Public API call methods

func (c *Client) Connect() (string, error) {
	// Instantiate login request & response struct
	options := LoginOptions(c.baseURL, c.apiKey)
	req, err := Login.GetRequest(options)
	if err != nil {
		return "", err
	}
	var resp = LoginResponse{}

	// Send login
	err = c.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}

	// Assign session key to client
	c.SessionKey = resp.Session
	// Log success
	slog.Debug(fmt.Sprintf("Succesful login, session key = %s", c.SessionKey))

	// Return session key
	return c.SessionKey, nil
}

func (c *Client) UploadFile(file string) (int, error) {
	// Create request
	options := UploadOptions(c.baseURL, c.SessionKey, file)
	req, err := Upload.GetRequest(options)
	if err != nil {
		return 0, err
	}
	resp := UploadResponse{}

	// Send upload request
	err = c.sendRequest(req, &resp)
	if err != nil {
		err = c.sessionCheck(err)
		if err != nil {
			return 0, err
		}
		return c.UploadFile(file)
	}

	// Return submission id
	return resp.SubID, nil
}

// Public methods that combine multiple calls

func (c *Client) GetPartialReview(subID int) (*PartialReview, error) {

	var partialReview = PartialReview{
		ID:       subID,
		Finished: false,
		Relevant: true,
	}

	// SubmissionStatus
	options := SubmissionStatusOptions(c.baseURL, subID)
	req, err := SubmissionStatus.GetRequest(options)
	if err != nil {
		return nil, err
	}

	resp := SubmissionStatusResponse{}
	if err = c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	if len(resp.JobCalibrations) > 0 {
		partialReview.Finished = true

		jobs := resp.Jobs
		// Job results request
		options := JobResultsOptions(c.baseURL, jobs[0])
		req, err := JobResults.GetRequest(options)
		if err != nil {
			return nil, err
		}

		resp := JobResultsResponse{}
		if err = c.sendRequest(req, &resp); err != nil {
			return nil, err
		}

		if resp.Status == "failure" {
			partialReview.Relevant = false
		} else {
			partialReview.Relevant = true
			partialReview.Calibration = resp.Calibration
		}
	}

	return &partialReview, nil
}

// Private utility

func (c *Client) checkSubmission(subID int) (*SubmissionStatusResponse, error) {

	// Create request
	options := SubmissionStatusOptions(c.baseURL, subID)
	req, err := SubmissionStatus.GetRequest(options)
	if err != nil {
		return nil, err
	}

	// Send request
	resp := SubmissionStatusResponse{}
	if err = c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Structs

type PartialReview struct {
	// Overview
	ID       int
	Finished bool
	Relevant bool
	// Tagged objects in field
	Objects []string
	// Telescope calibration
	Calibration CalibrationResponse
}

type FullReview struct {
	// Overview
	ID       int
	FileName string
	Finished bool
	Relevant bool
	// Start and finish times
	Start  string
	Finish string
	// Associated job IDs
	Jobs []int
	// Tags
	MachineTags    []string
	Tags           []string
	ObjectsInField []string
	// Object positions
	ObjectPositions AnnotationsResponse
	// Telescope calibration
	Calibration CalibrationResponse
}
