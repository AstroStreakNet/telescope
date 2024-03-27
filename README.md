# Telescope

Go wrapper around the astrometry API and image processing tools. Heavily focused around the FITS file format. 

```go

CropFITS()

```

## Astrometry API Interface

The Astrometry API can be interacted with through a 
client with inbuilt methods. This is designed to 
abstract away the majority of the actual API, simplifying
interaction to its most basic possible functions.

### Instancing

A client for the nova.astrometry.net API can be instanced
with the NewAstrometryClient function.

```go
var apiKey string = "example"
client := NewAstrometryClient(apiKey)
```

A client for a separately hosted astrometry instance can be 
instanced with the NewClient function.

```go
var baseURL string = "my.astrometry.instance"
var apiKey string = "example"
client := NewClient(baseURL, apiKey)
```

### Basic Interactions