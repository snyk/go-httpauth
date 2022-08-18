# go-httpauth

## Overview
This library introduces "advanced" HTTP Authentication mechanisms to be used in the golang HTTP stack (http.Transport). Current focus is on Proxy Authentication but future use is not limited to it.

The implementation supports automatic mechanism detection with `httpauth.AnyAuth`

Currently supported authentication mechanism: 
- Negotiate (see https://www.ietf.org/rfc/rfc4559.txt)
    - Kerberos (on all OS)
    - NTLM (on Windows)

## Usage
### Proxy Authentication
See `cmd/example1/main.go` 