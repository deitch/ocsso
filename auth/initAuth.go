package auth

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func InitAuth(server, clientVersion, userAgent string, client *http.Client) (string, string, string, string, string, error) {

	// create our vpn client
	if client == nil {
		client = &http.Client{}
	}

	initAuth, err := initAuthXML(server, clientVersion)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to marshal init-auth message: %v", err)
	}
	log.Debugf("init-auth request: %s", initAuth)

	initAuthReq, err := http.NewRequest("POST", server, bytes.NewReader(initAuth))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("could not create message to initialize authentication: %v", err)
	}
	// correct headers
	addAuthHeaders(initAuthReq, userAgent)

	resp, err := client.Do(initAuthReq)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error sending init-auth message: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("could not read init-auth body: %v", err)
	}
	log.Debugf("init-auth response: %s", body)
	initAuthResponse, err := parseAuthResponse(body)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("could not parse init-auth response: %v", err)
	}
	// extract the elements we need from the respData
	op := initAuthResponse.Opaque
	if op == nil {
		return "", "", "", "", "", errors.New("init-auth response missing <opaque> element")
	}
	tunnelGroup := op.TunnelGroup
	authMethod := op.AuthMethod
	configHash := op.ConfigHash
	log.Debugf("tunnel-group: %s", tunnelGroup)
	log.Debugf("auth-method: %s", authMethod)
	log.Debugf("config-hash: %s", configHash)

	authContent := initAuthResponse.Auth
	if authContent == nil {
		return "", "", "", "", "", errors.New("init-auth response missing <auth> element")
	}
	ssoLoginURL := authContent.SSOV2Login
	//ssoLoginFinal := authContent.SSOV2LoginFinal
	ssoCookieName := authContent.SSOV2TokenCookieName

	return tunnelGroup, authMethod, configHash, ssoLoginURL, ssoCookieName, nil
}

// initAuth create initial auth request
func initAuth(server, clientVersion string) *configAuth {
	version := version{
		Who:     "vpn",
		Version: clientVersion,
	}

	return &configAuth{
		Client:               "vpn",
		Type:                 "init",
		AggregateAuthVersion: 2,
		Version:              &version,
		DeviceId:             deviceID,
		GroupAccess:          server,
		Capabilities: &capabilities{
			AuthMethod: "single-sign-on-v2",
		},
	}
}

// initAuthXML create initial auth request as XML
func initAuthXML(server, clientVersion string) ([]byte, error) {
	ia := initAuth(server, clientVersion)
	initAuthMessage, err := xml.Marshal(ia)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal init-auth message: %v", err)
	}
	return initAuthMessage, nil
}
