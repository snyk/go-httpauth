package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/snyk/go-httpauth/pkg/httpauth"
	"github.com/snyk/go-httpauth/test/helper"
)

func runTest(proxyConfiguration func(req *http.Request) (*url.URL, error)) {
	logger := log.Default()

	// Create a new ProxyAuthenticator
	proxyAuthenticator := httpauth.NewProxyAuthenticator(httpauth.AnyAuth, proxyConfiguration, logger)
	transport := &http.Transport{
		DialContext: proxyAuthenticator.DialContext, // ensure to use the ProxyAuthenticator
		Proxy:       nil,                            // when using the ProxyAuthenticator, do not set the Transport Proxy
	}

	// create a client using the transport instance
	client := &http.Client{Transport: transport}

	// make a request
	response, err := client.Get("https://snyk.io")
	if err != nil {
		fmt.Println("Failed process request:", err)
	} else {
		fmt.Println("Received:", response.StatusCode)
	}
}

func main() {

	env := helper.NewProxyTestEnvironment("")

	if !env.HasDockerInstalled() {
		fmt.Println("This example requires docker to be installed.")
		os.Exit(1)
	} else {
		env.StartProxyEnvironment()

		os.Setenv("KRB5CCNAME", "FILE:"+env.CacheFile())
		os.Setenv("KRB5_CONFIG", env.ConfigFile())

		// this case depends on the environment variables HTTP_PROXY and HTTPS_PROXY, so it might not use the proxy if these are not set
		fmt.Println("--> Run with proxy Settings from environment Variables")
		runTest(http.ProxyFromEnvironment)

		// this case forces to use the proxy on localhost, which is started above
		fmt.Println("--> Run with proxy on localhost")
		runTest(func(req *http.Request) (*url.URL, error) { return url.Parse("http://localhost:3128") })

		// this case uses basic authentication
		fmt.Println("--> Run with basic proxy authentication from url ")
		runTest(func(req *http.Request) (*url.URL, error) { return url.Parse("http://user:password@localhost:3128") })

		env.StopProxyEnvironment()
	}
}
