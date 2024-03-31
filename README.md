# Telescope v0.1.3

Go wrapper around the astrometry API and image processing tools. Heavily focused around the FITS file format.

## FITS Manipulation

The telescope package can be utilised to accomplish some minor
manipulation of FITS data, such as cropping the data to a specific area
or converting it to a preferable format.

At the moment this is limited to converting a FITS file into a much smaller image format. 
The primary reason for doing this is displaying the data on a web page, creating
a 'copy' of the data in a JPEG or PNG format and using that instead of the FITS file
is vastly preferable.

```go
package main

import (
	"github.com/AstroStreakNet/telescope"
	"log"
)

func main() {
	// Must be .fits
	var inputFile = "file_to_convert.fits"
	// Output can be jpeg or png, simply provide a path of the desired type
	var outputPath = "file_converted.jpeg"

	err := telescope.ConvertFITS(inputFile, outputPath)
	if err != nil {
		log.Fatal(err)
	}
}
```

> Due to the nature of FITS files this image conversion will not be a perfect representation
of the actual data in the FITS file. 
> 
> This conversion only grabs the data in either the primary HDU or the first non-primary HDU. 
> For FITS files with multiple HDUs this will naturally result in an image that only has a select part of the overall data. 
> Due to this function being designed to help display FITS images on a website, only for public consumption, 
> this shortcoming was considered excusable.

For future iterations of this package it is hoped that the cropping of FITS files will be
implemented. Furthermore, the automatic cropping of the files to capture 'streaks' is the long-term
deliverable. Please message me with ideas.

## Astrometry API Interface

The Astrometry API can be interacted with through a 
client with inbuilt methods. This is designed to 
abstract away the majority of the actual API, simplifying
interaction to its most basic possible functions.

### Instancing a Client

A client for the nova.astrometry.net API can be instanced
with the **NewAstrometryClient()** function.

```go
package main

import (
	"github.com/AstroStreakNet/telescope/astrometry"
	"log"
)

func main() {
	var apiKey string = "KarshLovesWindowsAndGoogle"
	client := astrometry.NewAstrometryClient(apiKey)
}
```

If you are hosting your astrometry instance or have access to 
an alternatively hosted instance you can instead use the **NewClient()**
function.

```go
package main

import (
	"github.com/AstroStreakNet/telescope/astrometry"
	"log"
)

func main() {
	var baseURL string = "my.astrometry.instance/api"
	var apiKey string = "AdrianVsHR_CageFight"
	client := astrometry.NewClient(baseURL, apiKey)
}
```

### Uploading a file

```go
// Connect() returns session key
// Client automatically stores it, you don't have to use this return value
sessionKey, err := client.Connect()
if err != nil {
    log.Fatal(err)
}

// Uploading a file, can be any format supported by astrometry
// UploadFile() returns submission ID
// Client automatically stores it, don't have to use this return value
subID, err := c.UploadFile("./testdata/test_file.png")
if err != nil {
	log.Fatal(err)
}
```

> It is important to note that using the **Connect()** function is not strictly necessary.
> If you have no intention of uploading files or urls with the client then do not use this function.
>
> **Connect()** simply gets a session key, the astrometry session keys are only needed for uploading.

### Checking Submissions

Checking a submission can be done by providing the submission ID.
A struct with two fields is returned, allowing you to check if the job is done and
also query the job IDs.

```go
subStat, err := c.ReviewSubmission(subID)
if err != nil {
	log.Fatal(err)
}

subStat.Finished // bool
subStat.Jobs // []int list of job ids for submission
```
