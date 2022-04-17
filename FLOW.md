# SAML Process

1. Make a request to the server with the appropriate headers and body.
1. Read the info from the response.
1. Use the info from the response to open a Web page to the auth engine.
1. When the user is done logging in, read the appropriate response data from the browser.
1. Send the auth finish request to the URL provided.
1. Read the auth finish response, which has the credentials for openconnect.
1. Run openconnect with the provided credentials.

## Make a request to the server

The request is a `POST` request.

The server should be provided by the requestor, e.g. https://vpn.remote.com.

The headers and body must be specified correctly, or it will not work.

Headers:

```
User-Agent: "AnyConnect Linux_64 4.7.00136"
Accept: "*/*"
Accept-Encoding: "identity"
X-Transcend-Version: "1"
X-Aggregate-Auth: "1"
X-Support-HTTP-Auth: "true"
Content-Type: "application/x-www-form-urlencoded"
```

You can set the `User-Agent` to whatever you want, keeping in mind
that the server may be insistent that its client is AnyConnect.

The body must be xml, as follows:

```xml
<?xml version='1.0' encoding='UTF-8'?>
<config-auth client="vpn" type="init" aggregate-auth-version="2">
  <version who="vpn">4.7.00136</version>
  <device-id>linux-64</device-id>
  <group-select></group-select>
  <group-access>https://vpn.remote.com/</group-access>
  <capabilities>
    <auth-method>single-sign-on-v2</auth-method>
  </capabilities>
</config-auth>
```

## Read the info from the response

A successful response is `200`.

The response should be in XML format. An example follows:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<config-auth client="vpn" type="auth-request" aggregate-auth-version="2">
<opaque is-for="sg">
<tunnel-group>DefaultWEBVPNGroup</tunnel-group>
<auth-method>single-sign-on-v2</auth-method>
<config-hash>1646156124329</config-hash>
</opaque>
<auth id="main">
<title>Login</title>
<message>Please complete the authentication process in the AnyConnect Login window.</message>
<banner></banner>
<sso-v2-login>https://vpn.remote.com/+CSCOE+/saml/sp/login?tgname=DefaultWEBVPNGroup&#x26;acsamlcap=v2</sso-v2-login>
<sso-v2-login-final>https://vpn.remote.com/+CSCOE+/saml_ac_login.html</sso-v2-login-final>
<sso-v2-logout>https://vpn.remote.com/+CSCOE+/saml/sp/logout</sso-v2-logout>
<sso-v2-logout-final>https://vpn.remote.com/+CSCOE+/saml_ac_login.html</sso-v2-logout-final>
<sso-v2-token-cookie-name>acSamlv2Token</sso-v2-token-cookie-name>
<sso-v2-error-cookie-name>acSamlv2Error</sso-v2-error-cookie-name>
<form>
<input type="sso" name="sso-token"></input>
</form>
</auth>
</config-auth>
``` 

This contains everything we need to make the request, process the result, and submit the result to get credentials.

* `<tunnel-group>`: will be needed for final authentication
* `<auth-method>`: will be needed for final authentication
* `<config-hash>`: will be needed for final authentication
* `<sso-v2-login>`: where we should send the user to try to log in SSO
* `<sso-v2-login-final>`: where we should come back to, once we have logged into SSO and gotten the credentials we need
* `<sso-v2-token-cookie-name>`: the name of the cookie we need from the SSO login, that we will submit to `sso-v2-login-final`

There are others, which we will ignore for now.

## Open a Web page for authentication

Use the `sso-v2-login` URL to open a browser, and await completion.

## Read the response data

When the user is done, the browser window should be loaded with a cookie, whose name matches `sso-v2-cookie-name`. Read the contents of that cookie.

## Send the auth finish request

Send a `POST` request back to the original server URL. The headers and body must be correct.

Body:

```xml
<?xml version='1.0' encoding='UTF-8'?>
<config-auth client="vpn" type="auth-reply" aggregate-auth-version="2">
  <version who="vpn">4.7.00136</version>
  <device-id>linux-64</device-id>
  <session-token/>
  <session-id/>
  <opaque is-for="sg">
    <tunnel-group>DefaultWEBVPNGroup</tunnel-group>
    <auth-method>single-sign-on-v2</auth-method>
    <config-hash>1646156124329</config-hash>
  </opaque>
  <auth>
    <sso-token>71D2D28F0744FDEE74C74F6</sso-token>
  </auth>
</config-auth>
```

The contents of the important fields are:

* `<tunnel-group>`: matches what was returned from original request
* `<auth-method>`: matches what was returned from original request
* `<config-hash>`: matches what was returned from original request

Finally, the important one is the `<sso-token>`, whose contents are the value of the cookie.

The headers should be the same headers from the original request.

## Read the auth finish response

The response from the original server should be XML, return code `200`. For example:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<config-auth client="vpn" type="complete" aggregate-auth-version="2">
  <session-id>417792</session-id>
  <session-token>14B256@417792@4CF0@A0B9B161D548C12E5A01C01D736C36A6FF2F7833</session-token>
  <auth id="success">
    <banner>Content to show user</banner>
    <message id="0" param1="" param2=""></message>
  </auth>
  <capabilities>
    <crypto-supported>ssl-dhe</crypto-supported>
  </capabilities>
  <config client="vpn" type="private">
    <vpn-base-config>
      <nopkg></nopkg>
      <server-cert-hash>BD2CB3626D932E07026DBE6AB01394966F5D1BCE</server-cert-hash>
    </vpn-base-config>
    <opaque is-for="vpn-client"></opaque>
    <vpn-profile-manifest>
      <vpn rev="1.0">
        <file type="profile" service-type="user">
          <uri>/CACHE/stc/profiles/VPN.xml</uri>
          <hash type="sha1">44A3249DFA96589EB56496299813BAAB3B8BF572</hash>
        </file>
      </vpn>
    </vpn-profile-manifest>
    <vpn-language-manifest>
      <vpn rev="1.0">
        <file type="l10n" lang="en-us">
          <filename>AnyConnect.mo</filename>
          <hash type="sha1">503018D2008D95204B12939E48FC4133F5D58E4D</hash>
        </file>
      </vpn>
    </vpn-language-manifest>
  </config>
</config-auth>
```

There are a lot of fields here. The important ones:

* `config-auth/auth/banner`: contents should be displayed to the user.
* `config-auth/session-token`: contents are the token that will need to be used to authenticate to the vpn
* `config-auth/config/vpn-base-config/server-cert-hash`: contents are the server cert hash to be passed to openconnect

## Launch openconnect

Run openconnect:

```
openconnect --cookie <value-of-session-token> --servercert <value-of-server-cert-hash> https://vpn.remote.com/
```

For example:

```
openconnect --cookie 14B256@417792@4CF0@A0B9B161D548C12E5A01C01D736C36A6FF2F7833 --servercert BD2CB3626D932E07026DBE6AB01394966F5D1BCE https://vpn.remote.com/
```
