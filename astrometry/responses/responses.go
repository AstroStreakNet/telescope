package responses

// Error is naturally the structure of the Astrometry API's error messages
type Error struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errormessage"`
}

// Login provides the session ID, which is needed to upload files
type Login struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Session string `json:"session"`
}

// Upload includes the submission id, which is vital for monitoring the submissions progress
type Upload struct {
	Status string `json:"status"`
	SubID  int    `json:"subid"`
	Hash   string `json:"hash"`
}

// SubmissionStatus is the API's response when querying a submission's progress
type SubmissionStatus struct {
	ProcessingStarted  string  `json:"processing_started"`
	JobCalibrations    [][]int `json:"job_calibrations"`
	Jobs               []int   `json:"jobs"`
	ProcessingFinished string  `json:"processing_finished"`
	User               int     `json:"user"`
	UserImages         []int   `json:"user_images"`
}

// JobStatus is status of job, literally just one string, probably says 'success'
type JobStatus struct {
	Status string `json:"status"`
}

// Calibration represents the settings the telescope used to achieve the image
type Calibration struct {
	Parity      float64 `json:"parity"`
	Orientation float64 `json:"orientation"`
	PixelScale  float64 `json:"pixscale"`
	Radius      float64 `json:"radius"`
	Radian      float64 `json:"ra"`
	Decimal     float64 `json:"dec"`
}

// TaggedObjects is the response from getting job results on job objects
type TaggedObjects struct {
	Tags []string `json:"tags"`
}

// KnownObjects is an array of the known objects in the image, essentially same as taggedObjects
type KnownObjects struct {
	ObjectsInField []string `json:"objects_in_field"`
}

// Annotations is the known objects and their associated coordinates
type Annotations struct {
	List []struct {
		Radius float64  `json:"radius"`
		Type   string   `json:"type"`
		Names  []string `json:"names"`
		Pixelx float64  `json:"pixelx"`
		Pixely float64  `json:"pixely"`
	} `json:"annotations"`
}

// JobResults is all the resultant data from the API, excluding coordinates for objects
type JobResults struct {
	Status           string   `json:"status"`
	MachineTags      []string `json:"machine_tags"`
	Calibration      `json:"Calibration"`
	Tags             []string `json:"tags"`
	OriginalFilename string   `json:"original_filename"`
	ObjectsInField   []string `json:"objects_in_field"`
}
