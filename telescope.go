package telescope

// Data Preparation

// Crop image into desired dimensions
func Crop(image string, x, y, w, h float32) error {
	return nil
}

func CropAuto(image string) error {
	// TODO figure out how this would ever actually work
	return nil
}

func CropFITS(imagePath string) error {
	// TODO this will be very painful to do
	return nil
}

// ConvertFITS is a function for taking a FITS file and outputting the data in a different image format, jpeg/png.
// Takes path to FITS file, creates new file in jpeg/png format at designated output path
func ConvertFITS(imagePath, outputPath string) error {
	return nil
}

// Data conversion/preparation
// At some point or another it may be necessary to create functions for manipulating the image data.
// For example:
// - Function which will convert the FITS data, lossless, to an equivalent format.
// - Function to normalise FITS data for ml purposes
// - Function to flatten FITS data for ml purposes
