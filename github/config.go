package github

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/logging"
	"golang.org/x/oauth2"
)

type Config struct {
	Token        string
	Organization string
	BaseURL      string
	Insecure     bool
}

type Organization struct {
	name        string
	client      *github.Client
	StopContext context.Context
}

// Client configures and returns a fully initialized GithubClient
func (c *Config) Client() (interface{}, error) {
	var org Organization
	org.name = c.Organization
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)

	ctx := context.Background()

	if c.Insecure {
		httpClient := insecureHttpClient()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	}

	tc := oauth2.NewClient(ctx, ts)

	tc.Transport = logging.NewTransport("Github", tc.Transport)
	tc.Transport = antiAbuseTransport(1*time.Second, tc.Transport)

	org.client = github.NewClient(tc)
	if c.BaseURL != "" {
		u, err := url.Parse(c.BaseURL)
		if err != nil {
			return nil, err
		}
		org.client.BaseURL = u
	}

	return &org, nil
}

func insecureHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// This is to prevent triggering abuse rate limits
// https://developer.github.com/v3/guides/best-practices-for-integrators/#dealing-with-abuse-rate-limits
func antiAbuseTransport(sleep time.Duration, originalTransport http.RoundTripper) http.RoundTripper {
	// TODO: mutex
	// TODO: sleep 1 sec
}
