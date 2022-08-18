package httpauth

import (
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
	mechanism := StringFromAuthenticationMechanism(a.activeMechanism)

	if a.activeMechanism == Negotiate { // supporting mechanism: Negotiate (SPNEGO)
		var token string
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
			mechanisms, _ := GetMechanismsFromHttpFieldValue(authorizeValue)
			a.logger.Printf("Authorization to %s using: %s", url, mechanisms)
		}
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
