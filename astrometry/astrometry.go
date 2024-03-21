package astrometry

// Constants

const baseURL string = "http://nova.astrometry.net/api"
const signInURL string = baseURL + "/login"
const uploadURL string = baseURL + "/upload"

// Abstractions

type Client struct {
	sessionKey  string
	submissions []string
}

func Connect(apiKey string) (Client, error) {
	sessionKey, err := signIn(apiKey)
	if err != nil {
		return Client{}, err
	}
	return Client{
		sessionKey: sessionKey,
	}, nil
}

func (c *Client) UploadFile(file string) {

}

// Data structs

// Calibration represents the settings the telescope used to achieve the image
type Calibration struct {
	parity      float32
	orientation float32
	pixelScale  float32
	radius      float32
	radians     float32
	degrees     float32
}

// KnownObject represents an object found within the analysed image
type KnownObject struct {
	radius   float32
	category string
	names    []string
	x        float32
	y        float32
}

type Results struct {
	status      string
	tags        []string
	calibration Calibration
	fileName    string
	objects     []string
}

// API functions

// signIn
func signIn(apiKey string) (string, error) {

	return "", nil
}

// upload
func upload(sessionKey, file string) {

}

// submissionStatus
func submissionStatus(sessionKey, subID string) {

}

// jobStatus
func jobStatus(sessionKey, jobID string) {

}

// jobResultsCalibration
func jobResultsCalibration(sessionKey, jobID string) {

}

// jobResultsTags
func jobResultsTags(sessionKey, jobID string) []string {
	return []string{"2", "3"}
}

// jobResultsObjects
func jobResultsObjects(sessionKey, jobID string) []string {
	return nil
}

// jobResultsObjectsCoordinates
func jobResultsObjectsCoordinates(sessionKey, jobID string) []KnownObject {
	return nil
}

// jobResults
func jobResults(sessionKey, jobID string) {

}
