package httpauth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/snyk/error-catalog-golang-public/cli"
	"github.com/snyk/error-catalog-golang-public/snyk_errors"
	"github.com/stretchr/testify/assert"
)

// closedAddress returns an address that is guaranteed to refuse connections.
func closedAddress(t *testing.T) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)
	address := listener.Addr().String()
	assert.Nil(t, listener.Close())
	return address
}

func expectSingleAuthorizationAttempt(handler *MockAuthenticationHandlerInterface, proxyUrl *url.URL) {
	handler.EXPECT().SetLogger(testLogger).Times(1)
	handler.EXPECT().IsStopped().Return(false).Times(1)
	handler.EXPECT().GetAuthorizationValue(proxyUrl, "").Return("", nil).Times(1)
	handler.EXPECT().Cancel().Times(1)
	handler.EXPECT().IsStopped().Return(true).Times(1)
	handler.EXPECT().Close().Times(1)
}

func assertProxyConnectionError(t *testing.T, err error) snyk_errors.Error {
	t.Helper()

	var snykError snyk_errors.Error
	assert.True(t, errors.As(err, &snykError))
	assert.Equal(t, cli.NewProxyConnectionError("").ErrorCode, snykError.ErrorCode)
	return snykError
}

func Test_ProxyAuthenticator_DialContext_unreachableProxy_isProxyConnectionError(t *testing.T) {
	proxyUrl := &url.URL{Scheme: "http", User: url.UserPassword("someone", "s3cr3t"), Host: closedAddress(t)}
	proxyFunc := func(request *http.Request) (*url.URL, error) { return proxyUrl, nil }

	ctrl := gomock.NewController(t)
	mockedAuthenticationHandler := NewMockAuthenticationHandlerInterface(ctrl)
	expectSingleAuthorizationAttempt(mockedAuthenticationHandler, proxyUrl)

	authenticator := NewProxyAuthenticator(Negotiate, proxyFunc, testLogger)
	authenticator.CreateHandler = func(mechanism AuthenticationMechanism) AuthenticationHandlerInterface {
		return mockedAuthenticationHandler
	}

	connection, err := authenticator.DialContext(context.Background(), "tcp", "snyk.io:443")

	assert.Nil(t, connection)

	snykError := assertProxyConnectionError(t, err)
	assert.Contains(t, snykError.Detail, "http://someone:xxxxx@"+proxyUrl.Host)
	assert.NotContains(t, snykError.Detail, "s3cr3t")

	// the underlying network error stays available to errors.As
	var opError *net.OpError
	assert.True(t, errors.As(err, &opError))
}

func Test_ProxyAuthenticator_DialContext_unresolvableProxy_isProxyConnectionError(t *testing.T) {
	cause := fmt.Errorf("invalid proxy configuration")
	proxyFunc := func(request *http.Request) (*url.URL, error) { return nil, cause }

	authenticator := NewProxyAuthenticator(Negotiate, proxyFunc, testLogger)

	connection, err := authenticator.DialContext(context.Background(), "tcp", "snyk.io:443")

	assert.Nil(t, connection)

	snykError := assertProxyConnectionError(t, err)
	assert.Equal(t, "proxy connection failed: invalid proxy configuration", snykError.Detail)
	assert.True(t, errors.Is(err, cause))
}

func Test_ProxyAuthenticator_DialContext_withoutProxy_isNotProxyConnectionError(t *testing.T) {
	proxyFunc := func(request *http.Request) (*url.URL, error) { return nil, nil }

	authenticator := NewProxyAuthenticator(Negotiate, proxyFunc, testLogger)

	connection, err := authenticator.DialContext(context.Background(), "tcp", closedAddress(t))

	assert.Nil(t, connection)
	assert.NotNil(t, err)

	var snykError snyk_errors.Error
	assert.False(t, errors.As(err, &snykError))
}

func Test_NewProxyConnectionError(t *testing.T) {
	cause := fmt.Errorf("connection refused")

	t.Run("mentions the proxy it failed to connect through", func(t *testing.T) {
		err := NewProxyConnectionError(&url.URL{Scheme: "http", Host: "proxy.internal:8080"}, cause)

		snykError := assertProxyConnectionError(t, err)
		assert.Equal(t, "proxy connection to http://proxy.internal:8080 failed: connection refused", snykError.Detail)
		assert.True(t, errors.Is(err, cause))
	})

	t.Run("handles an unknown proxy", func(t *testing.T) {
		err := NewProxyConnectionError(nil, cause)

		snykError := assertProxyConnectionError(t, err)
		assert.Equal(t, "proxy connection failed: connection refused", snykError.Detail)
	})
}
