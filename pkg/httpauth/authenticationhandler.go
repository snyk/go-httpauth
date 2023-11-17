package httpauth

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
)

type AuthenticationHandlerInterface interface {
	Close()
	Cancel()
	Succesful()
	IsStopped() bool
	GetAuthorizationValue(url *url.URL, responseToken string) (string, error)
	Update(availableMechanism map[AuthenticationMechanism]string) (string, error)
	SetBasicAuthentication(userInfo *url.Userinfo) error
	SetLogger(logger *log.Logger)
	SetSpnegoProvider(spnegoProvider SpnegoProvider)
}

type AuthenticationHandler struct {
	spnegoProvider  SpnegoProvider
	Mechanism       AuthenticationMechanism
	activeMechanism AuthenticationMechanism
	state           AuthenticationState
	cycleCount      int
	logger          *log.Logger
	userInfo        *url.Userinfo
}

func NewHandler(mechanism AuthenticationMechanism) AuthenticationHandlerInterface {
	a := &AuthenticationHandler{
		spnegoProvider:  SpnegoProviderInstance(),
		Mechanism:       mechanism,
		activeMechanism: mechanism,
		state:           Initial,
		logger:          log.New(io.Discard, "", 0),
	}

	return a
}

func (a *AuthenticationHandler) Close() {
	a.spnegoProvider.Close()
	a.state = Close
}

func (a *AuthenticationHandler) GetAuthorizationValue(url *url.URL, responseToken string) (authorizeValue string, err error) {
	var token string

	if url.User != nil {
		// As soon as a user and password are given inside the URL, we assume that basic authentication is being required.
		err = a.SetBasicAuthentication(url.User)
		if err != nil {
			a.logger.Println(err)

			// in case of anyauth, supress the error to try other authentication mechanisms
			if a.activeMechanism == AnyAuth {
				err = nil
			}
		}
	}

	mechanism := StringFromAuthenticationMechanism(a.activeMechanism)
	if a.activeMechanism == Negotiate { // supporting mechanism: Negotiate (SPNEGO)
		var done bool

		if len(responseToken) == 0 && Negotiating == a.state {
			a.state = Error
			return "", fmt.Errorf("Authentication failed! Unexpected empty token during negotiation!")
		}

		a.state = Negotiating

		token, done, err = a.spnegoProvider.GetToken(url, responseToken)
		if err != nil {
			a.state = Error
			return "", err
		}

		if done {
			a.logger.Println("Local security context established!")
		}

		authorizeValue = mechanism + " " + token

		if len(token) > 0 {
			mechanisms, errMechanism := GetMechanismsFromHttpFieldValue(authorizeValue)
			a.logger.Printf("Authorization to %s using: %s (err=%v)", url, mechanisms, errMechanism)
		}
	} else if a.activeMechanism == BasicAuth { // supporting mechanism: Basic
		password, _ := a.userInfo.Password()
		userPass := fmt.Sprintf("%s:%s", a.userInfo.Username(), password)
		token = base64.StdEncoding.EncodeToString([]byte(userPass))
		authorizeValue = mechanism + " " + token

		a.logger.Printf("Authorization to %s using: [%s]", url, mechanism)
	}

	a.cycleCount++
	if a.cycleCount >= maxCycleCount {
		err = fmt.Errorf("Failed to authenticate within %d cycles, stopping now!", maxCycleCount)
		authorizeValue = ""
	}

	return authorizeValue, err
}

func (a *AuthenticationHandler) Update(availableMechanism map[AuthenticationMechanism]string) (responseToken string, err error) {

	// if AnyAuth is selected, we need to determine the best supported mechanism on both sides
	if a.activeMechanism == AnyAuth {
		// currently we only support Negotiate, AnyAuth will use Negotiate if the communication partner proposes it
		if _, ok := availableMechanism[Negotiate]; ok {
			a.activeMechanism = Negotiate
			a.logger.Printf("Selected Mechanism: %s\n", StringFromAuthenticationMechanism(a.activeMechanism))
		}
	}

	// extract the token for the active mechanism
	if token, ok := availableMechanism[a.activeMechanism]; ok {
		responseToken = token
	} else {
		err = fmt.Errorf("Incorrect or unsupported Mechanism detected! %s", availableMechanism)
	}

	return responseToken, err
}

func (a *AuthenticationHandler) SetBasicAuthentication(userInfo *url.Userinfo) error {
	if _, hasPassword := userInfo.Password(); hasPassword {
		a.userInfo = userInfo
		a.activeMechanism = BasicAuth
		return nil
	} else {
		return fmt.Errorf("Failed to set basic authentication, missing password")
	}
}

func (a *AuthenticationHandler) IsStopped() bool {
	return (a.state == Done || a.state == Error || a.state == Cancel || a.state == Close)
}

func (a *AuthenticationHandler) Cancel() {
	a.state = Cancel
	a.logger.Println("AuthenticationHandler.Cancel()")
}

func (a *AuthenticationHandler) Succesful() {
	a.state = Done
	a.logger.Println("AuthenticationHandler.Succesful()")
}

func (a *AuthenticationHandler) SetLogger(logger *log.Logger) {
	a.logger = logger
	a.spnegoProvider.SetLogger(logger)
}

func (a *AuthenticationHandler) SetSpnegoProvider(spnegoProvider SpnegoProvider) {
	a.spnegoProvider = spnegoProvider
}
