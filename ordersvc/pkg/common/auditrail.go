package common

import (
	"bytes"
	"io"
	"net/http"
	"slices"
)

type AudiTrail struct {
	URI       string
	Method    string
	Payload   string
	IpAddress string
	UserAgent string
}

// NewAudiTrail will create a new AudiTrail object
func NewAudiTrail(req *http.Request, realIP string) AudiTrail {
	var payload string

	// List of sensitive endpoints that should not be logged
	sensitiveEndpoints := []string{
		"/users/login",
	}

	// Check if current URI contains any sensitive endpoints
	if !slices.Contains(sensitiveEndpoints, req.RequestURI) {
		reqBody, err := io.ReadAll(req.Body)
		if err == nil {
			// Cleanup the payload string by removing \n\t and extra spaces
			cleanPayload := bytes.ReplaceAll(reqBody, []byte("\n"), []byte(""))
			cleanPayload = bytes.ReplaceAll(cleanPayload, []byte("\t"), []byte(""))
			cleanPayload = bytes.Join(bytes.Fields(cleanPayload), []byte(" "))
			payload = string(cleanPayload)

			// Create a new reader from the body bytes and restore it to req.Body
			// This is needed because reading the body consumes it
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
	}

	audiTrail := AudiTrail{
		URI:       req.RequestURI,
		Method:    req.Method,
		IpAddress: realIP,
		UserAgent: req.UserAgent(),
		Payload:   payload,
	}

	return audiTrail
}
