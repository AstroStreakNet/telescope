package telescope

import (
	"errors"
	"github.com/AstroStreakNet/telescope/util/FITS"
	"github.com/astrogo/fitsio"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

// Data Preparation

// Crop image into desired dimensions
func Crop(image string, x, y, w, h float32) error {
	return nil
}

func CropAuto(image string) (string, error) {
	// TODO figure out how this would ever actually work
	return "", nil
}

func CropFITS(imagePath string) error {
	// TODO this will be very painful to do
	return nil
}

// ConvertFITS is a function for taking a FITS file and outputting the data in a different image format, jpeg/png.
// Takes path to FITS file, creates new file in jpeg/png format at designated output path
func ConvertFITS(inputPath, outputPath string) error {

	fits, err := FITS.OpenFits(inputPath)
	if err != nil {
		return err
	}
	defer func(fits *fitsio.File) {
		err := fits.Close()
		if err != nil {
			log.Fatalf("ReadCloser failure: %s", err)
		}
	}(fits)

	var fImg fitsio.Image
	var ok bool

	if len(fits.HDUs()) > 1 {
		fImg, ok = fits.HDU(1).(fitsio.Image)
	} else {
		fImg, ok = fits.HDU(0).(fitsio.Image)
	}

	if !ok {
		return errors.New("failure to read FITS data as image")
	}
	img := fImg.Image()
	// Check header first before doing this.
	// Flip image, due to way FITS stores data.
	// img = imaging.FlipV(img)

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	switch filepath.Ext(outputPath) {
	case ".png":
		err = png.Encode(f, img)
	case ".jpeg":
		err = jpeg.Encode(f, img, nil)
	case ".fits":
		return errors.New("fits files don't need to be converted to fits files, fool")
	default:
		return errors.New("unsupported format")
	}
	if err != nil {
		return err
	}
	return nil
}

// Data conversion/preparation
// At some point or another it may be necessary to create functions for manipulating the image data.
// For example:
// - Function which will convert the FITS data, lossless, to an equivalent format.
// - Fpack/Funpack FITS compression
// - Function to normalise FITS data for ml purposes
// - Function to flatten FITS data for ml purposes

// Data extraction

type HeaderData struct {
	// Coordinate
	ObsID string
	RA    string
	DEC   string
	MJD   float64
	// Image data
	Radius       float64
	ExposureTime float64
}

func GetHeaderData(fitsFile string) (*HeaderData, error) {
	fits, err := FITS.OpenFits(fitsFile)
	if err != nil {
		return nil, err
	}
	defer func(fits *fitsio.File) {
		err := fits.Close()
		if err != nil {
			log.Fatalf("ReadCloser failure: %s", err)
		}
	}(fits)

	// Header
	header := fits.HDU(0).Header()
	data := new(HeaderData)
	// Query header for keywords
	data.ObsID = FITS.HeaderQuery("OBSID", header).(string)
	data.RA = FITS.HeaderQuery("RA", header).(string)
	data.DEC = FITS.HeaderQuery("DEC", header).(string)
	data.MJD = FITS.HeaderQuery("MJD-OBS", header).(float64)
	data.Radius = FITS.HeaderQuery("RADIUS", header).(float64)
	data.ExposureTime = FITS.HeaderQuery("EXPTIME", header).(float64)

	return data, nil
}
