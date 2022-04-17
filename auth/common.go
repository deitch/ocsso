package auth

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	deviceID = "linux-64"
)

// parseAuthResponse parse XML to get a finish-auth-response
func parseAuthResponse(b []byte) (*configAuth, error) {
	auth := configAuth{}
	if err := xml.Unmarshal(b, &auth); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %v", err)
	}
	return &auth, nil
}

func addAuthHeaders(req *http.Request, userAgent string) {
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "identity")
	req.Header.Add("X-Transcend-Version", "1")
	req.Header.Add("X-Aggregate-Auth", "1")
	req.Header.Add("X-Support-HTTP-Auth", "true")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
}
