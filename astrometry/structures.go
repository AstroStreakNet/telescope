package astrometry

// ErrorResponse is naturally the structure of the Astrometry API's error messages
type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errormessage"`
}

// LoginResponse provides the session ID, which is needed to upload files
type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Session string `json:"session"`
}

// UploadResponse includes the submission id, which is vital for monitoring the submissions progress
type UploadResponse struct {
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

// JobStatus is status of job, literally just one string
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

// KnownObjects is an array of the known objects in the image, essentially same as TaggedObjects
type KnownObjects struct {
	ObjectsInField []string `json:"objects_in_field"`
}

// Annotation is the position data of a tagged object in the image
type Annotation struct {
	Radius float64  `json:"radius"`
	Type   string   `json:"type"`
	Names  []string `json:"names"`
	PixelX float64  `json:"pixelx"`
	PixelY float64  `json:"pixely"`
}

// AnnotationsList is the known objects and their associated coordinates
type AnnotationsList struct {
	Annotations []Annotation `json:"annotations"`
}

// JobResults is all the resultant data from the API, excluding coordinates for objects
type JobResults struct {
	Status           string      `json:"status"`
	MachineTags      []string    `json:"machine_tags"`
	Calibration      Calibration `json:"calibration"`
	Tags             []string    `json:"tags"`
	OriginalFilename string      `json:"original_filename"`
	ObjectsInField   []string    `json:"objects_in_field"`
}
