package astrometry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Endpoint URLS
// These are separated from the base URL for testing purposes. Also, if the API is ever hosted locally the changeover
// will be significantly quicker

const (
	loginEP         = "/login"                    // Sign in to API, returns session key
	uploadEP        = "/upload"                   // Upload file to api, returns submission ID
	subStatusEP     = "/submissions/%s"           // Get status of submission, returns job IDs
	jobStatusEP     = "/jobs/%s"                  // Get status of specific job
	calibrationEP   = "/jobs/%s/calibration"      // Get calibration results of specific job
	taggedObjectsEP = "/jobs/%s/machine_tags"     // Get tags for picture
	knownObjectsEP  = "/jobs/%s/objects_in_field" // Get known objects in field of picture, essentially just tags
	annotationsEP   = "/jobs/%s/annotations"      // Get known objects and their positions
	jobResultsEP    = "/jobs/%s/info"             // Get results of specific job

)

// Request Structures
// If the endpoint requires arguments, then they are here

const (
	loginR = `request-json={"apikey": "%s"}`
	// All the options are naturally put as no, even if the user has consented to our site they have not consented to
	// their images being on, or used by, Astrometry
	uploadR = `request-json={"session": "%s", "allow_commercial_use": "n", "allow_modifications": "n", "publicly_visible": "n"}`
)

// doTheStuff is a Client method that sends off requests, checks if the response is an error message,
// and if not then decodes the data into a specified data structure
func (c *Client) doTheStuff(request *http.Request, dataStructure interface{}) error {
	// Send provided request
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	// If no error in sending request than defer close
	defer resp.Body.Close()

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
	err = json.NewDecoder(&respBytes).Decode(&dataStructure)
	if err != nil {
		return err
	}

	// Return nothing if no errors have occurred throughout whole process
	return nil
}

// Endpoints

// login is a POST request to the API, providing the API key, which allows for a session key to be granted
func (c *Client) login() (*LoginResponse, error) {

	// Define arguments for request
	arguments := []byte(fmt.Sprintf(loginR, c.apiKey))
	// Define request structure
	req, err := http.NewRequest("POST", c.BaseURL+loginEP, bytes.NewBuffer(arguments))
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)
	// Set necessary headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response = LoginResponse{}
	// Send request, decode data into LoginResponse structure
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// upload is a function for posting an image file to the Astrometry API
func (c *Client) upload(filePath string) (*UploadResponse, error) {

	// Open file to get data
	fileData, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileData.Close()

	// Get file name from the provided path
	fileName := filepath.Base(filePath)

	// Create multipart form from file data
	// Create multipart writer
	var reqBody bytes.Buffer
	w := multipart.NewWriter(&reqBody)

	err = w.WriteField("request-json", fmt.Sprintf(uploadR, c.SessionKey))
	if err != nil {
		return nil, err
	}

	fileBuffer, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	// Copy data from file to multipart form body
	_, err = io.Copy(fileBuffer, fileData)
	if err != nil {
		return nil, err
	}

	w.Close()

	// Create request
	req, err := http.NewRequest(
		"POST",
		c.BaseURL+uploadEP,
		&reqBody,
	)
	if err != nil {
		return nil, err
	}

	// Set request headers
	req.Header.Set("Content-Type", w.FormDataContentType())

	var response = UploadResponse{}
	// Send request, decode data into UploadResponse structure
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getSubmissionStatus queries the Astrometry API for the current status of a submission
func (c *Client) getSubmissionStatus(subID string) (*SubmissionStatus, error) {

	// Define request, submission ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(subStatusEP, subID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = SubmissionStatus{}
	// Send request, decode data into SubmissionStatus structure
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getJobStatus, returns singular string struct that will most likely just say "success"
func (c *Client) getJobStatus(jobID string) (*JobStatus, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(jobStatusEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = JobStatus{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getCalibration, returns calibration data of telescope used to take photo
func (c *Client) getCalibration(jobID string) (*Calibration, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(calibrationEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = Calibration{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getTaggedObjects, returns objects that have been tagged in the image
func (c *Client) getTaggedObjects(jobID string) (*TaggedObjects, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(taggedObjectsEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = TaggedObjects{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getKnownObjects, essentially the same as getTaggedObjects. I'm unsure as to why the Astrometry API has decided to
// define two endpoints that are basically the same
func (c *Client) getKnownObjects(jobID string) (*KnownObjects, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(knownObjectsEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = KnownObjects{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getAnnotations
func (c *Client) getAnnotations(jobID string) (*AnnotationsList, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(annotationsEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = AnnotationsList{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getJobResults
func (c *Client) getJobResults(jobID string) (*JobResults, error) {

	// Define request, job ID has to be included in URL
	req, err := http.NewRequest(
		"GET",
		c.BaseURL+fmt.Sprintf(jobResultsEP, jobID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// req.WithContext(ctx)

	var response = JobResults{}
	err = c.doTheStuff(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
