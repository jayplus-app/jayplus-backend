package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetSubdomain(r *http.Request) (string, error) {
	domainParts := strings.Split(r.Host, ".")

	switch {
	case r.Host == "localhost":
		return "", errors.New("No subdomain found")
	case len(domainParts) == 2 && domainParts[1] == "localhost":
		return domainParts[0], nil
	case len(domainParts) < 3:
		return "", errors.New("No subdomain found")
	case domainParts[0] == "www" && len(domainParts) == 3:
		return "", errors.New("No subdomain found")
	case domainParts[0] == "www":
		return domainParts[1], nil
	default:
		return domainParts[0], nil
	}
}
