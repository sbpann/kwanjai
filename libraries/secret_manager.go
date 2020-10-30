package libraries

import (
	"fmt"
	"kwanjai/config"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// AccessSecretVersion function returns secret value (string) and error.
func AccessSecretVersion(name string) (string, error) {
	InitializeGCP()
	client, err := secretmanager.NewClient(config.Context)
	if err != nil {
		return "error", fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(config.Context, req)
	if err != nil {
		return "error", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
