package astrometry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AstroStreakNet/telescope/astrometry/endpoints"
	"github.com/AstroStreakNet/telescope/astrometry/responses"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client

type Client struct {
	apiKey      string
	baseURL     string
	SessionKey  string
	submissions map[string]int
	finished    map[string]int
	httpClient  *http.Client
}

// Client factory functions

func NewAstrometryClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "http://nova.astrometry.net/api", // Astrometry doesn't support https
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		submissions: make(map[string]int),
		finished:    make(map[string]int),
	}
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		submissions: make(map[string]int),
		finished:    make(map[string]int),
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
	var errorResponse = responses.Error{}
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
	req, err := endpoints.Login.Request(c.baseURL, c.apiKey, "")
	if err != nil {
		return "", err
	}
	var resp = responses.Login{}

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
	req, err := endpoints.UploadFile.Request(c.baseURL, c.SessionKey, file)
	if err != nil {
		return 0, err
	}
	resp := responses.Upload{}

	// Send upload request
	err = c.sendRequest(req, &resp)
	if err != nil {
		err = c.sessionCheck(err)
		if err != nil {
			return 0, err
		}
		return c.UploadFile(file)
	}

	// If successful, add submission id to client list
	c.addSubmission(file, resp.SubID)
	// Return submission id
	return resp.SubID, nil
}

// Public methods that combine multiple calls

func (c *Client) GetPartialReview(file string) (*PartialReview, error) {

	subID := c.getSubmissionID(file)

	var partialReview = PartialReview{
		ID:       subID,
		FileName: file,
		Finished: false,
		Relevant: true,
	}

	// SubmissionStatus
	req, err := endpoints.SubmissionStatus.Request(c.baseURL, strconv.Itoa(subID), "")
	if err != nil {
		return nil, err
	}

	resp := responses.SubmissionStatus{}
	if err = c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	if len(resp.JobCalibrations) > 0 {
		partialReview.Finished = true

		jobs := resp.Jobs
		// Calibration request
		req, err := endpoints.Calibration.Request(c.baseURL, strconv.Itoa(jobs[0]), "")
		if err != nil {
			return nil, err
		}

		resp := responses.Calibration{}
		if err = c.sendRequest(req, &resp); err != nil {
			return nil, err
		}

		partialReview.Calibration = resp
	}

	return &partialReview, nil
}

func (c *Client) CheckSubmission(file string) (*responses.SubmissionStatus, error) {
	// Check if submission is listed
	subID := c.getSubmissionID(file)
	if subID == 0 {
		return nil, errors.New("submission not in client submission/finished list")
	}

	// Create request
	req, err := endpoints.SubmissionStatus.Request(c.baseURL, strconv.Itoa(subID), "")
	if err != nil {
		return nil, err
	}

	// Send request
	resp := responses.SubmissionStatus{}
	if err = c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	// Check if submission has finished
	if len(resp.JobCalibrations) > 0 {
		c.removeSubmission(file)
		c.addFinished(file, subID)
	}

	return &resp, nil
}

// Client utility functions, boilerplate getters and setters

func (c *Client) UpdateAllSubmissions() error {
	for file := range c.CurrentSubmissions() {
		if _, err := c.CheckSubmission(file); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) CurrentSubmissions() map[string]int {
	return c.submissions
}

func (c *Client) FinishedSubmissions() map[string]int {
	return c.finished
}

func (c *Client) getSubmissionID(fileName string) int {
	var subID int
	subID = c.submissions[fileName]
	if subID == 0 {
		subID = c.finished[fileName]
	}
	return subID
}

func (c *Client) addSubmission(fileName string, subID int) {
	c.submissions[fileName] = subID
}

func (c *Client) removeSubmission(fileName string) {
	delete(c.submissions, fileName)
}

func (c *Client) addFinished(fileName string, subID int) {
	c.finished[fileName] = subID
}

func (c *Client) removeFinished(fileName string) {
	delete(c.finished, fileName)
}

// Structs

type PartialReview struct {
	// Overview
	ID       int
	FileName string
	Finished bool
	Relevant bool
	// Tagged objects in field
	Objects []string
	// Telescope calibration
	Calibration responses.Calibration
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
	ObjectPositions responses.Annotations
	// Telescope calibration
	Calibration responses.Calibration
}
