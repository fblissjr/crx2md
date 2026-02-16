package crx

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const downloadURLTemplate = "https://clients2.google.com/service/update2/crx?response=redirect&prodversion=131.0&acceptformat=crx2,crx3&x=id%%3D%s%%26installsource%%3Dondemand%%26uc"

// extensionIDPattern matches a 32-character lowercase alphanumeric Chrome extension ID.
var extensionIDPattern = regexp.MustCompile(`^[a-z]{32}$`)

// webStoreURLPattern extracts extension ID from Chrome Web Store URLs.
var webStoreURLPattern = regexp.MustCompile(`chromewebstore\.google\.com/detail/[^/]+/([a-z]{32})`)

// ParseExtensionID extracts the extension ID from various input formats:
// - Chrome Web Store URL
// - Raw extension ID
func ParseExtensionID(input string) (string, error) {
	// Direct extension ID
	if extensionIDPattern.MatchString(input) {
		return input, nil
	}

	// Chrome Web Store URL
	if strings.Contains(input, "chromewebstore.google.com") {
		matches := webStoreURLPattern.FindStringSubmatch(input)
		if len(matches) >= 2 {
			return matches[1], nil
		}
		return "", fmt.Errorf("could not extract extension ID from URL: %s", input)
	}

	return "", fmt.Errorf("not a valid extension ID or Chrome Web Store URL: %s", input)
}

// Download fetches a CRX file from the Chrome Web Store and saves it to a temp file.
// Returns the path to the temp file.
func Download(extensionID string) (string, error) {
	url := fmt.Sprintf(downloadURLTemplate, extensionID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download extension: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d (extension ID: %s)", resp.StatusCode, extensionID)
	}

	tmpFile, err := os.CreateTemp("", "crx2md-*.crx")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("write temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
