package FITS

import (
	"errors"
	"github.com/astrogo/fitsio"
	"image"
	"io"
)

// TODO somehow modify FITS data in Go

func Decode(reader io.Reader) (image.Image, error) {

	fits, err := fitsio.Open(reader)
	if err != nil {
		return nil, err
	}

	var fitsImage fitsio.Image
	var ok bool

	if len(fits.HDUs()) > 1 {
		fitsImage, ok = fits.HDU(1).(fitsio.Image)
	} else {
		fitsImage, ok = fits.HDU(0).(fitsio.Image)
	}

	if !ok {
		return nil, errors.New("failure to read FITS data as image")
	}

	return fitsImage.Image(), nil
}
