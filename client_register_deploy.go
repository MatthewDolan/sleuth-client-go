package sleuth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (client *Client) RegisterDeploy(
	ctx context.Context,
	deploymentSlug string,
	environment string,
	sha string,
	deployedAt time.Time,
) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(client.newRegisterDeployRequest(
		true,
		environment,
		sha,
		deployedAt,
	)); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.registerDeployURL(deploymentSlug),
		buf,
	)
	if err != nil {
		return err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	return nil
}

func (client *Client) registerDeployURL(deploymentSlug string) string {
	return fmt.Sprintf("https://app.sleuth.io/api/1/%s/%s/register_deploy", client.organizationSlug, deploymentSlug)
}

// https://help.sleuth.io/sleuth-api#rest-api-details
type registerDeployRequest struct {
	IgnoreIfDuplicate string `json:"ignore_if_duplicate"`
	Tags              string `json:"tags"`
	Environment       string `json:"environment"`
	Email             string `json:"email"`
	Date              string `json:"date"`
	APIKey            string `json:"api_key"`
	SHA               string `json:"sha"`
}

func (client *Client) newRegisterDeployRequest(
	ignoreIfDup bool,
	environment string,
	sha string,
	deployedAt time.Time,
) *registerDeployRequest {
	// TODO add email & tag representations

	req := &registerDeployRequest{
		APIKey:      client.apiKey,
		Environment: environment,
		// ISO 8601 deployment date string
		// https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
		// https://stackoverflow.com/questions/522251/whats-the-difference-between-iso-8601-and-rfc-3339-date-formats
		Date: deployedAt.Format(time.RFC3339),
		SHA:  sha,
	}
	if ignoreIfDup {
		// If the value is provided and set to "true" Sleuth won't return a 400
		// if we see a SHA that has already been registered
		req.IgnoreIfDuplicate = "true"
	}
	return req
}
