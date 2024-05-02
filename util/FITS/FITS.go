package FITS

import (
	"errors"
	"github.com/astrogo/fitsio"
	"log"
	"os"
	"path/filepath"
)

// TODO somehow modify FITS data in Go

func OpenFits(path string) (*fitsio.File, error) {
	if filepath.Ext(path) != ".fits" {
		return nil, errors.New("provided file is not a FITS file, is not .fits")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("ReadCloser failure: %s", err)
		}
	}(file)

	return fitsio.Open(file)
}

func HeaderQuery(key string, header *fitsio.Header) any {
	value := header.Get(key)
	if value != nil {
		return value.Value
	}
	return nil
}
