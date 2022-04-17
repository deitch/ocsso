package auth

import (
	"encoding/xml"
)

type configAuth struct {
	XMLName              xml.Name      `xml:"config-auth"`
	Client               string        `xml:"client,attr"`
	Type                 string        `xml:"type,attr"`
	AggregateAuthVersion int           `xml:"aggregate-auth-version,attr"`
	Version              *version      `xml:"version,omitempty"`
	DeviceId             string        `xml:"device-id,omitempty"`
	GroupSelect          string        `xml:"group-select,omitempty"`
	GroupAccess          string        `xml:"group-access,omitempty"`
	Capabilities         *capabilities `xml:"capabilities,omitempty"`
	Opaque               *opaque       `xml:"opaque,omitempty"`
	Auth                 *auth         `xml:"auth,omitempty"`
	SessionID            string        `xml:"session-id,omitempty"`
	SessionToken         string        `xml:"session-token,omitempty"`
	Config               *config       `xml:"config,omitempty"`
}

type version struct {
	XMLName xml.Name `xml:"version"`
	Who     string   `xml:"who,attr,omitempty"`
	Version string   `xml:",chardata"`
}
type capabilities struct {
	XMLName         xml.Name `xml:"capabilities"`
	AuthMethod      string   `xml:"auth-method"`
	CryptoSupported string   `xml:"crypto-supported,omitempty"`
}

type opaque struct {
	XMLName     xml.Name `xml:"opaque"`
	IsFor       string   `xml:"is-for,attr"`
	TunnelGroup string   `xml:"tunnel-group"`
	AuthMethod  string   `xml:"auth-method"`
	ConfigHash  string   `xml:"config-hash"`
}

type auth struct {
	XMLName xml.Name `xml:"auth"`
	ID      string   `xml:"id,attr,omitempty"`
	Title   string   `xml:"title,omitempty"`
	Message string   `xml:"message,omitempty"`
	Banner  string   `xml:"banner,omitempty"`

	SSOV2Login       string `xml:"sso-v2-login,omitempty"`
	SSOV2LoginFinal  string `xml:"sso-v2-login-final,omitempty"`
	SSOV2Logout      string `xml:"sso-v2-logout,omitempty"`
	SSOV2LogoutFinal string `xml:"sso-v2-logout-final,omitempty"`

	SSOV2TokenCookieName string `xml:"sso-v2-token-cookie-name,omitempty"`
	SSOV2ErrorCookieName string `xml:"sso-v2-error-cookie-name,omitempty"`

	SSOToken string `xml:"sso-token,omitempty"`
	Form     *form  `xml:"form,omitempty"`
}

type form struct {
	XMLName xml.Name `xml:"form"`
	Input   input    `xml:"input"`
}
type input struct {
	XMLName xml.Name `xml:"input"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:",chardata"`
}

type config struct {
	XMLName             xml.Name             `xml:"config"`
	Client              string               `xml:"client,attr"`
	Type                string               `xml:"type,attr"`
	VPNBaseConfig       *vpnBaseConfig       `xml:"vpn-base-config,omitempty"`
	VPNProfileManifest  *vpnProfileManifest  `xml:"vpn-profile-manifest,omitempty"`
	VPNLanguageManifest *vpnLanguageManifest `xml:"vpn-language-manifest,omitempty"`
	Opaque              *opaque              `xml:"opaque,omitempty"`
}
type vpnBaseConfig struct {
	XMLName        xml.Name `xml:"vpn-base-config"`
	NoPkg          string   `xml:"nopkg,omitempty"`
	ServerCertHash string   `xml:"server-cert-hash,omitempty"`
}
type vpnProfileManifest struct {
	XMLName xml.Name `xml:"vpn-profile-manifest"`
	VPN     *vpn     `xml:"vpn,omitempty"`
}
type vpn struct {
	XMLName xml.Name `xml:"vpn"`
	File    *file    `xml:"file"`
}
type file struct {
	XMLName     xml.Name `xml:"file"`
	Type        string   `xml:"type,attr"`
	ServiceType string   `xml:"service-type,attr"`
	URI         string   `xml:"uri"`
	Hash        *hash    `xml:"hash"`
}
type hash struct {
	XMLName xml.Name `xml:"hash"`
	Type    string   `xml:"type,attr"`
	Content string   `xml:",chardata"`
}
type vpnLanguageManifest struct {
	XMLName xml.Name `xml:"vpn-language-manifest"`
	VPN     *vpn     `xml:"vpn,omitempty"`
}
