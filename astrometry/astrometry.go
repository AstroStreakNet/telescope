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
	"telescope/astrometry/endpoints"
	"telescope/astrometry/responses"
	"time"
)

// Constants

const actualAPI = "https://nova.astrometry.com/api"

// Client

type Client struct {
	apiKey      string
	baseURL     string
	SessionKey  string
	submissions []int
	finished    []int
	httpClient  *http.Client
}

func (c *Client) addSubmission(s int) {
	c.submissions = append(c.submissions, s)
}

func NewAstrometryClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: actualAPI,
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

func (c *Client) reconnect(e error) (bool, error) {
	// Check if session key has expired.
	if strings.Contains(e.Error(), "no session with key") {
		// If session key expired, login again
		_, err := c.Connect()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Public methods

func (c *Client) Connect() (string, error) {
	// Instantiate login request & response struct
	req, err := endpoints.Login.Request(c.baseURL, c.apiKey)
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
	resp := responses.Upload{}
	// Send upload request
	err = c.sendRequest(req, &resp)
	if err != nil {

		strings.Contains(err.Error(), "no session with key")

	}
	// If successful, add submission id to client list
	c.addSubmission(resp.SubID)
	// Return submission id
	return resp.SubID, nil
}

func (c *Client) ReviewFile(subID string) {

}

func (c *Client) CurrentSubmissions() []int {
	return c.submissions
}

func (c *Client) FinishedSubmissions() []int {
	return c.finished
}
