package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func SetDomainEnv() {
	frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL"))
	if frontendURL == "" {
		return
	}

	u, err := url.Parse(frontendURL)
	if err != nil || u.Host == "" {
		fmt.Fprintf(os.Stderr, "Invalid FRONTEND_URL env variable: %s\n Expected format: http(s)://<host>[:port] \n > Example: https://structured-notes.nl \n > Example: http://localhost:8200\n", frontendURL)
		os.Exit(1)
	}

	// If there is a port, remove it for cookie domain
	host := u.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	_ = os.Setenv("DOMAIN_CLIENT", fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	_ = os.Setenv("COOKIE_DOMAIN", host)
}
