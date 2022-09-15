package httpauth

import (
	"strings"
)

type AuthenticationMechanism string
type AuthenticationState int

const maxCycleCount int = 10

const (
	NoAuth           AuthenticationMechanism = "noauth"
	Negotiate        AuthenticationMechanism = "negotiate"
	AnyAuth          AuthenticationMechanism = "anyauth"
	BasicAuth        AuthenticationMechanism = "basic"
	UnknownMechanism AuthenticationMechanism = "unknownmechanism"
)

const (
	Initial     AuthenticationState = iota
	Negotiating AuthenticationState = iota
	Done        AuthenticationState = iota
	Error       AuthenticationState = iota
	Cancel      AuthenticationState = iota
	Close       AuthenticationState = iota
)

const (
	AuthorizationKey      string = "Authorization"
	ProxyAuthorizationKey string = "Proxy-Authorization"
	ProxyAuthenticateKey  string = "Proxy-Authenticate"
)

func StringFromAuthenticationMechanism(mechanism AuthenticationMechanism) string {
	return strings.Title(string(mechanism))
}

func AuthenticationMechanismFromString(mechanism string) AuthenticationMechanism {
	tmp := strings.ToLower(mechanism)
	return AuthenticationMechanism(tmp)
}

func GetMechanismAndToken(HttpFieldValue string) (AuthenticationMechanism, string) {
	mechanism := UnknownMechanism
	token := ""

	authenticateValue := strings.Split(HttpFieldValue, " ")
	if len(authenticateValue) >= 1 {
		mechanism = AuthenticationMechanismFromString(authenticateValue[0])
	}

	if len(authenticateValue) == 2 {
		token = authenticateValue[1]
	}

	return mechanism, token
}

func IsSupportedMechanism(mechanism AuthenticationMechanism) bool {
	if mechanism == Negotiate || mechanism == AnyAuth {
		return true
	}
	return false
}
