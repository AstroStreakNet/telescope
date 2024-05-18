package astrometry

// Response structures with JSON decoding tags

type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errormessage"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Session string `json:"session"`
}

type UploadResponse struct {
	Status string `json:"status"`
	SubID  int    `json:"subid"`
	Hash   string `json:"hash"`
}

type SubmissionStatusResponse struct {
	ProcessingStarted  string  `json:"processing_started"`
	JobCalibrations    [][]int `json:"job_calibrations"`
	Jobs               []int   `json:"jobs"`
	ProcessingFinished string  `json:"processing_finished"`
	User               int     `json:"user"`
	UserImages         []int   `json:"user_images"`
}

type CalibrationResponse struct {
	Parity      float64 `json:"parity"`
	Orientation float64 `json:"orientation"`
	PixelScale  float64 `json:"pixscale"`
	Radius      float64 `json:"radius"`
	Radian      float64 `json:"ra"`
	Decimal     float64 `json:"dec"`
}

type AnnotationsResponse struct {
	Annotations []struct {
		Radius float64  `json:"radius"`
		Type   string   `json:"type"`
		Names  []string `json:"names"`
		PixelX float64  `json:"pixelx"`
		PixelY float64  `json:"pixely"`
	} `json:"annotations"`
}

type JobResultsResponse struct {
	Status           string              `json:"status"`
	MachineTags      []string            `json:"machine_tags"`
	Calibration      CalibrationResponse `json:"calibration"`
	Tags             []string            `json:"tags"`
	OriginalFilename string              `json:"original_filename"`
	ObjectsInField   []string            `json:"objects_in_field"`
}
