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
	submissions []int
	finished    []int
	httpClient  *http.Client
}

// Client utility functions

func (c *Client) addSubmission(s int) {
	c.submissions = append(c.submissions, s)
}

// removeSubmission is much slower than the alternative removeSubmissionByIndex
func (c *Client) removeSubmission(s int) {
	for i, subID := range c.submissions {
		if subID == s {
			c.removeSubmissionByIndex(i)
		}
	}
}

func (c *Client) removeSubmissionByIndex(index int) {
	c.submissions = append(c.submissions[:index], c.submissions[index+1:]...)
}

func (c *Client) addFinished(s int) {
	c.finished = append(c.finished, s)
}

func (c *Client) removeFinished(s int) {
	for i, subID := range c.finished {
		if subID == s {
			c.removeFinishedByIndex(i)
		}
	}
}

func (c *Client) removeFinishedByIndex(index int) {
	c.finished = append(c.finished[:index], c.finished[index+1:]...)
}

// Client factory functions

func NewAstrometryClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "http://nova.astrometry.net/api",
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

// Public methods

func (c *Client) Connect() (string, error) {
	// Instantiate login request & response struct
	req, err := endpoints.Login.Request(c.baseURL, c.apiKey)
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
	req, err := endpoints.UploadFile.Request(c.baseURL, c.SessionKey)
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
	c.addSubmission(resp.SubID)
	// Return submission id
	return resp.SubID, nil
}

func (c *Client) ReviewSubmission(subID int) (*SubStat, error) {
	req, err := endpoints.SubmissionStatus.Request(c.baseURL, strconv.Itoa(subID))
	if err != nil {
		return nil, err
	}

	resp := responses.SubmissionStatus{}
	if err = c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	subStat := SubStatFromResponse(resp)
	if subStat.Finished {
		c.removeSubmission(subID)
		c.addFinished(subID)
	}

	return SubStatFromResponse(resp), nil
}

func (c *Client) CheckSubmissions() ([]int, error) {
	var finished []int
	for i, cSub := range c.CurrentSubmissions() {
		subStat, err := c.ReviewSubmission(cSub)
		if err != nil {
			return finished, err
		}
		if subStat.Finished {
			c.removeSubmissionByIndex(i)
			c.addFinished(cSub)
			finished = append(finished, cSub)
		}
	}
	return finished, nil
}

func (c *Client) CurrentSubmissions() []int {
	return c.submissions
}

func (c *Client) FinishedSubmissions() []int {
	return c.finished
}

// Structs

type SubStat struct {
	Finished bool
	Jobs     []int
}

func SubStatFromResponse(resp responses.SubmissionStatus) *SubStat {
	subStat := SubStat{
		Finished: false,
		Jobs:     resp.Jobs,
	}
	// If job calibrations not empty then job has been finished
	if len(resp.JobCalibrations) > 0 {
		subStat.Finished = true
	}
	return &subStat
}
