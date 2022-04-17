package auth

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func FinishAuth(server, clientVersion, tunnelGroup, authMethod, configHash, token, userAgent string, client *http.Client) (string, string, string, error) {
	var banner string
	finishAuth, err := finishAuthXML(server, clientVersion, tunnelGroup, authMethod, configHash, token)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal finish-auth message: %v", err)
	}
	log.Debugf("finish-auth request: %s", finishAuth)

	finishAuthReq, err := http.NewRequest("POST", server, bytes.NewReader(finishAuth))
	if err != nil {
		return "", "", "", fmt.Errorf("could not create message to finish authentication: %v", err)
	}
	// correct headers
	addAuthHeaders(finishAuthReq, userAgent)

	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(finishAuthReq)
	if err != nil {
		return "", "", "", fmt.Errorf("error sending finish-auth message: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("could not read finish-auth body: %v", err)
	}
	log.Debugf("finish-auth response: %s", body)

	finishAuthResponse, err := parseAuthResponse(body)
	if err != nil {
		return "", "", "", fmt.Errorf("could not unmarshal: %v", err)
	}

	if finishAuthResponse.SessionToken == "" {
		return "", "", "", fmt.Errorf("finish-auth-response session-token is empty")
	}
	if finishAuthResponse.Config == nil {
		return "", "", "", fmt.Errorf("finish-auth-response did not include final config")
	}
	if finishAuthResponse.Config.VPNBaseConfig == nil {
		return "", "", "", fmt.Errorf("finish-auth-response did not include config>vpn-base-config")
	}
	if finishAuthResponse.Config.VPNBaseConfig.ServerCertHash == "" {
		return "", "", "", fmt.Errorf("finish-auth-response config>vpn-base-config>server-cert-hash is empty")
	}
	if finishAuthResponse.Auth != nil {
		banner = finishAuthResponse.Auth.Banner
	}
	return finishAuthResponse.SessionToken, finishAuthResponse.Config.VPNBaseConfig.ServerCertHash, banner, nil
}

// finishAuth create finish auth request
func finishAuth(server, clientVersion, tunnelGroup, authMethod, configHash, token string) *configAuth {
	version := version{
		Who:     "vpn",
		Version: clientVersion,
	}

	return &configAuth{
		Client:               "vpn",
		Type:                 "auth-reply",
		AggregateAuthVersion: 2,
		Version:              &version,
		DeviceId:             deviceID,
		Opaque: &opaque{
			IsFor:       "sg",
			TunnelGroup: tunnelGroup,
			AuthMethod:  authMethod,
			ConfigHash:  configHash,
		},
		Auth: &auth{
			SSOToken: token,
		},
	}
}

// finishAuthXML create finish auth request as XML
func finishAuthXML(server, clientVersion, tunnelGroup, authMethod, configHash, token string) ([]byte, error) {
	fa := finishAuth(server, clientVersion, tunnelGroup, authMethod, configHash, token)
	finishAuthMessage, err := xml.Marshal(fa)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal finish-auth message: %v", err)
	}
	return finishAuthMessage, nil
}
