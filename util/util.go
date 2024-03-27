package util

import (
	"os"
)

// Testing utilities

// GetTestDataString is primarily used to write the json data into a http response for client request testing
func GetTestDataString(filePath string) string {
	data, err := os.ReadFile("../testdata" + filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// GetTestData is used to load the json data of the response, used to verify client's returned data is correct
func GetTestData(filePath string) []byte {
	data, err := os.ReadFile("../testdata" + filePath)
	if err != nil {
		panic(err)
	}
	return data
}

func GetTestFilePath(filePath string) string {
	return "../testdata" + filePath
}
