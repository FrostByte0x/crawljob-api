package handler

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"strings"
)

func ValidateUrl(requestedUrl string) error {
	// load list of valid domains to download files
	allowedDomains := strings.Split(os.Getenv("ALLOWED_DOMAINS"), ",")
	if allowedDomains == nil {
		slog.Warn("ALLOWED_DOMAINS is empty, all downloads will be rejected!")
	}
	// Parse the URL
	parsedUrl, err := url.Parse(requestedUrl)
	if err != nil {
		return err
	}
	if parsedUrl.Scheme != "https" {
		return fmt.Errorf("Requested URL cannot be parsed as an https host")
	}
	if !slices.Contains(allowedDomains, parsedUrl.Host) {
		return fmt.Errorf("Requested Host is not in the allowed domains. Contact the server owner to allow the domain!")
	}
	return nil
}
