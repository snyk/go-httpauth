package httpauth

import (
	"fmt"
	"net/url"

	"github.com/snyk/error-catalog-golang-public/cli"
	"github.com/snyk/error-catalog-golang-public/snyk_errors"
)

// ProxyConnectionError indicates that a connection failed while resolving, reaching or
// authenticating against the configured proxy, rather than while reaching the target host.
func NewProxyConnectionError(proxyUrl *url.URL, err error) error {
	detail := fmt.Sprintf("proxy connection failed: %v", err)
	if proxyUrl != nil {
		detail = fmt.Sprintf("proxy connection to %s failed: %v", proxyUrl.Redacted(), err)
	}

	return cli.NewProxyConnectionError(detail, snyk_errors.WithCause(err))
}
