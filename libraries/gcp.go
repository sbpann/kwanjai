package libraries

import (
	"kwanjai/config"
	"os"
)

// InitializeGCP from credential.json.
func InitializeGCP() {
	if os.Getenv("GIN_MODE") == "release" {
		return
	}
	defaultCredential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if defaultCredential == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.BaseDirectory+"/.secret/credential.json")
	}
}
