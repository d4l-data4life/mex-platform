# Pairgen


**WARNING: Never deploy the Pairgen service! It is only to be used during development and testing!**

## Background

Requests to the MEx services APIs must be authenticated by a JWT.
Each JWT must be signed by a trusted public key.
This trust is established by injecting one or more public keys (or public key certificates) into the service.
This injection cannot happen at deployment time because we are authenticating against Microsoft Azure Active Directory (AAD) and they rotate the keys in regular intervals.
Microsoft AAD exposes a URL under which the keys can be retrieved.
The above-mentioned trust hence extends to the key retrieval URL: it must be HTTPS and the server certificate must be signed by a known/trusted certificate authority (CA).
The configuration value `MEX_JWKS_REMOTE_KEYS_URI` contains this key retrieval URL and the MEx services enforce it to be HTTPS and the server certificate to be valid.

## Problem

This is all fine for the Microsoft AAD keys endpoint (https://login.microsoftonline.com/common/discovery/keys), but becomes a problem for integration tests.
Integration tests spawn MEx services and execute API requests.
In order for these requests to be accepted, they must carry a valid JWT.

Three options come to mind:

- Option 1: Add a configuration toggle to the MEx services to disable JWT validation.
  - Verdict: This is clearly untenable because it would introduce a single bit of information on which the entire security hinges.
- Option 2: Use client credentials grant flow to obtain a proper JWT from Microsoft AAD which can be validated with keys obtained from the above keys endpoint.
  - Verdict: This requires that the integration test code have access to a OAuth client secret which is required during the client credentials grant flow. This confidential information cannot (well, at least, should not) be baked into the integration test code and must be injected at test run time. This makes the setup for the tests harder as the client secret must be injected via Vault or similar.
- Option 3: Execute a service which acts as a JWT provider just for the integration tests. This is what the Pairgen service is for!

## Service description

During startup the Pairgen service will generate one RSA key pair (hence the name (_pair_ _gen_)erator) for each key ID specified in the environment variable `KEY_IDS`.
The public keys of these pairs can be retrieved in JWKS format via the endpoint `/public_keys` (this corresponds to the Microsoft AAD keys endpoint above).
However, this `/public_keys` endpoint must be HTTPS with a verifiable CA signature to be accepted by the MEx services.

This is why Pairgen during startup also synthesizes two further key pairs:

- One which is used to generate a custom CA certificate.
- Another one that is used as TLS certificate to serve the `/public_keys` endpoint via HTTPS.
  - This certificate will be signed by the custom CA certificate.

In order for the MEx services to verify the TLS certificate of the `/public_keys` endpoint, it must be told about the custom CA certificate.
The Pairgen service can be configured to export the custom CA certificate as a PEM file (`CA_CERT_FILE`) and the MEx services can be configured to load additional CA certificates from a PEM file (`MEX_WEB_ADDITIONAL_CA_CERT_FILE`).
This way the custom CA certificate can be shared in a joint Docker volume during integration testing.

## Service endpoints

Pairgen will open two ports (`HTTP_PORT` and `HTTPS_PORT`, defaulting to 3000 and 3001, respectively) under which the following endpoints are available.

### HTTPS

#### GET `/public_keys`

This is the only HTTPS endpoint.
It returns the public keys in JWKS format.

Example:

If the service was started with `KEY_IDS='foo bar'` (note, there is always a `default` key pair):

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "foo",
      "use": "sig",
      "alg": "RS256",
      "e": "AQAB",
      "n": "y3dmQVRnXPFvr4STeAdZQJA1sqMQTYF4s_gwBffWbmwLf2uM3349Dd-iHBmVQu9ew34OGlT6vcuKl4zyGOS-is0A_6pLRWXTVDog9FW-HRw8faNy9woHR1liBff3c0M_pXafCV836xk60uXMSK6iHEmTgbtWomj2OjvmG1j8oCE"
    },
    {
      "kty": "RSA",
      "kid": "bar",
      "use": "sig",
      "alg": "RS256",
      "e": "AQAB",
      "n": "sVmoVrcJe390uyclk832vQ_vnhws3lOrY-j_h4ny68DNxxUTn9J-NizxDKJx2gYPmyOaF6hgtZEHjT2vORMoETNlAkZwfa-Gj5ImdUDCjlHoFZzvxM4gFHuCRH4fylm1sYw3yZ4DikVo_WxafX2XHLhhaDXFoWeTEusEQ4-ixt0"
    },
    {
      "kty": "RSA",
      "kid": "default",
      "use": "sig",
      "alg": "RS256",
      "e": "AQAB",
      "n": "0EU-EFmNQ5yB1OSzuMVtgRG4l1OEKlrBUWP2-zidbc_xpeLG5IOCL2cOneHC7szS5YJaq4rHLa7L_OS3_ifHM0TbG_mgWQ0DYzm9ZY6WX3zzRQ1Jta5rSGwdUgtXINR7ixKb569W9OFgjxAJkdL0ym8SqXjS4ubj8ILt2suj5ak"
    }
  ]
}
```

### HTTP

#### GET `/`

This route serves as a readiness endpoint.
If it returns a HTTP 200, the Pairgen service is ready to be used.

---

#### GET `/ca`

In addition to writing the custom CA certificate out to the file `CA_CERT_FILE`, it can be retrived by this endpoint:

Example:

```
-----BEGIN CERTIFICATE-----
MIIDuDCCAqCgAwIBAgIBATANBgkqhkiG9w0BAQsFADBoMQ8wDQYDVQQDEwZNRXgg
Q0ExCzAJBgNVBAYTAkRFMRQwEgYDVQQIEwtCcmFuZGVuYnVyZzEQMA4GA1UEBxMH
UG90c2RhbTEPMA0GA1UEChMGTUV4IENBMQ8wDQYDVQQLEwZNRXggQ0EwHhcNMjIw
MzE3MTMwNzM4WhcNMjMwMzE3MTMwNzM4WjBoMQ8wDQYDVQQDEwZNRXggQ0ExCzAJ
BgNVBAYTAkRFMRQwEgYDVQQIEwtCcmFuZGVuYnVyZzEQMA4GA1UEBxMHUG90c2Rh
bTEPMA0GA1UEChMGTUV4IENBMQ8wDQYDVQQLEwZNRXggQ0EwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQCXz3tq3aYheP/MkMVrsy9I9nnRG8ID+lnUYsE1
d9ZrAc1+E48hLUCLNBgjf2YlLOTblUTBy5k3uy58OAn1w39r9cCNM98DQJOofnYj
v1cviSbMx6JwB0le33o9/1Axc9IVOOXJC/ZElb/MHa/VjO22nSIDCB1twItKAiYN
6XluFIlYA4Pv8lDKljBhIpIvqacdd1VJ5JqNxXVSJB7dDPH9DaHkoExaJtKsENlc
uuAlatkTyNydMlSrJ1pVOOegVjt7AqXS8ovbnKZIj9m1hWX+4yDr/EqoqoHwUVIx
fbCWcxjjurxukuI0j0EjFYnmUAC1H5lF/Fc9UtQCNHEdbjVlAgMBAAGjbTBrMAwG
A1UdEwQFMAMBAf8wCwYDVR0PBAQDAgL0MDsGA1UdJQQ0MDIGCCsGAQUFBwMBBggr
BgEFBQcDAgYIKwYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcDCDARBglghkgBhvhC
AQEEBAMCAPcwDQYJKoZIhvcNAQELBQADggEBACe2CEu0j72O8ggQbhX0q/4sN/NO
GnigrNus4LC4G3+XJMQtB37MxzeGgPuxyky1vLuI92O+YHtSvWMc2OhSvp7IJKr4
w9T8WJihNPaX/KWX1x7/9WGm4/1PxSiLLOdzueB59LM86eOvUzUpa6eah/cWsbh5
nNYU1i12sx5xwxj3tu3ClU2gqRACjQ2/h3byJF7eI9WfRF9nSWkikw1ln5WA+mSh
N6Bmr18ErSg+QyJqh3ogL7RqGJjz96mhk8GLnY0tgSkde3xG6sJtr0KObzENh7e4
Qzsnb1j9eW1PIm81wvuLTEyQB2g3XFkNbt021PDJ0KDkLGRZHaf748h6ROw=
-----END CERTIFICATE-----
```

---

#### GET `/private_keys/:kid`

This endpoint returns the private key of the key pair with ID `kid` in PEM format.

---

#### GET `/public_keys/:kid`

This endpoint returns the public key of the key pair with ID `kid` in PEM format.

---

#### POST `/jwts/:kid`

This endpoint takes the POST body as claims and returns a JWT signed with the private key of key pair `kid`.

Example:

```sh
curl -i -X POST -d '{"hello":"world"}' http://localhost:5080/jwts/foo

HTTP/1.1 200 OK
X-Powered-By: Express
Content-Type: text/plain; charset=utf-8
X-Auth-Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImZvbyJ9.eyJpYXQiOjE2NDc1MjI5MzR9.gbDiwEW4J2a1doh9r2KOoutKMooEus612HlajMekWd2r7Sqp3vhK7evFlg2019GjR_yGMUJftEKwma4HjLYw7rGvX0GcQRBvbTg3UK3MpF9bJOBX98Bo0LkEdhMQ_68aAyFy0Q8w1ePiSQCfGJDcl-x19Bj-0p4fnrcsP3ocHvw
Content-Length: 249
ETag: W/"f9-3y1yYU8M0ABRUrHb7HareC1poTo"
Date: Thu, 17 Mar 2022 13:15:34 GMT
Connection: keep-alive

eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImZvbyJ9.eyJpYXQiOjE2NDc1MjI5MzR9.gbDiwEW4J2a1doh9r2KOoutKMooEus612HlajMekWd2r7Sqp3vhK7evFlg2019GjR_yGMUJftEKwma4HjLYw7rGvX0GcQRBvbTg3UK3MpF9bJOBX98Bo0LkEdhMQ_68aAyFy0Q8w1ePiSQCfGJDcl-x19Bj-0p4fnrcsP3ocHvw
```

The JWT is returned in the response body as well as in the `X-Auth-Token` header.
The latter locationi simplifies using the endpoint with tools like [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) which allows to reference parts of a previous response to use in subsequent requests.

---

#### POST `/token`

This endpoint mimics an OAuth client credentials grant flow.
The `client_id` serves as the key ID (`kid`).

Example:

```sh
curl -i -X POST -d 'client_id=foo&oid=me' http://localhost:5080/token

HTTP/1.1 200 OK
X-Powered-By: Express
Content-Type: application/json; charset=utf-8
Content-Length: 307
ETag: W/"133-j2GLvcKKRexN2HyyHMtpKyTQeJw"
Date: Thu, 17 Mar 2022 13:26:11 GMT
Connection: keep-alive

{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImZvbyJ9.eyJjbGllbnRfaWQiOiJmb28iLCJvaWQiOiJtZSIsImlhdCI6MTY0NzUyMzU3MX0.VQDlR_kZh5ldd-Gcy7IHJgZuRKJWBqi9_8pZX3-Tb9Zv7tvQFH0WO84fj3yb4RXxt0XjoYLMQru-PtuGCb1dyRenKAymuPq61Tm7jD8Lj_2VRaEks86xzl4m3OTIltGm9_IpLOvzi4C227qmVfV9I-x12ArmB0agDzEGxS06DEQ"}
```

---

## Security aspects

- The MEx services always verify JWTs of incoming requests.
- The MEx services always enforce the `MEX_JWKS_REMOTE_KEYS_URI` to be HTTPS.
- The MEx services always validate the server certificate of `MEX_JWKS_REMOTE_KEYS_URI`.
- There is no configuration toggle to bypass any of this.
- **Never deploy the Pairgen service!** It is only to be used during development and testing!
