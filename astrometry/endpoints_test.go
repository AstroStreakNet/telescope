package astrometry

import (
	"testing"
)

func TestEndpointGetRequest(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
	}{
		{"login", Login},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testEnd := testSetup()
			defer testEnd()

		})
	}
}
