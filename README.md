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

## Contributing

To ensure the long-term stability and quality of this project, we are moving to a closed-contribution model effective August 2025. This change allows our core team to focus on a centralized development roadmap and rigorous quality assurance, which is essential for a component with such extensive usage.

All of our development will remain public for transparency. We thank the community for its support and valuable contributions.

## Getting Support

GitHub issues have been disabled on this repository as part of our move to a closed-contribution model. The Snyk support team does not actively monitor GitHub issues on any Snyk development project.

For help with the Snyk CLI or Snyk in general, please use the [Snyk support page](https://support.snyk.io/), which is the fastest way to get assistance.