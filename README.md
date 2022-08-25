# go-httpauth

## Overview
This library introduces "advanced" HTTP Authentication mechanisms to be used in the golang HTTP stack (http.Transport). While the golang HTTP stack provides support for Basic Authentication or more generally authentication mechanisms that require a single authentication message to be send. This implementation adds support for mechanisms that require multiple messages to be exchanged for authentication, like challenge response based types like NTLM.

The current focus is on Proxy Authentication but future use is not limited to it.

The implementation supports automatic mechanism detection with `httpauth.AnyAuth`

Currently supported authentication mechanism: 
- Negotiate (see https://www.ietf.org/rfc/rfc4559.txt)
    - Kerberos (on all OS)
    - NTLM (on Windows)

## Usage
### Proxy Authentication
See `cmd/example1/main.go` 