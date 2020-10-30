package libraries

import (
	"kwanjai/config"
	"os"
)

var isInitialzed bool = false

// InitializeGCP from credential.json.
func InitializeGCP() {
	if isInitialzed {
		return
	}
	defaultCredential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if defaultCredential == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.BaseDirectory+"/.secret/credential.json")
	}
	isInitialzed = true
}
