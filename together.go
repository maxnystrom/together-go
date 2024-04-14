package together

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultScheme   = "https"
	defaultHostname = "api.together.xyz"
	defaultBasePath = "/"
	userAgent       = "together-go"

	errEmptyAPIToken = "invalid credentials: API Token must not be empty" //nolint:gosec,unused
)

var (
	Version string = "v1"
)

type API struct {
	APIKey    string
	BaseURL   string
	UserAgent string
	headers   http.Header
	Client    *retryablehttp.Client
	Debug     bool
}

func New(key string) (*API, error) {
	if key == "" {
		return nil, errors.New(errEmptyAPIToken)
	}

	api := &API{}
	api.BaseURL = fmt.Sprintf("%s://%s", defaultScheme, defaultHostname)
	api.Client = retryablehttp.NewClient()
	api.Client.RetryMax = 5
	api.Debug = false
	api.APIKey = key
	api.UserAgent = userAgent + "/" + Version

	return api, nil
}

func (api *API) request(ctx context.Context, method, uri string, reqBody io.Reader, headers http.Header) (*http.Response, error) {
	req, err := retryablehttp.NewRequestWithContext(ctx, method, api.BaseURL+uri, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed: %w", err)
	}

	combinedHeaders := make(http.Header)
	copyHeader(combinedHeaders, api.headers)
	copyHeader(combinedHeaders, headers)
	req.Header = combinedHeaders

	req.Header.Set("Authorization", "Bearer "+api.APIKey)

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if api.Debug {
		dump, err := httputil.DumpRequestOut(req.Request, true)
		if err != nil {
			return nil, err
		}

		// Strip out any sensitive information from the request payload.
		sensitiveKeys := []string{api.APIKey}
		for _, key := range sensitiveKeys {
			if key != "" {
				valueRegex := regexp.MustCompile(fmt.Sprintf("(?m)%s", key))
				dump = valueRegex.ReplaceAll(dump, []byte("[redacted]"))
			}
		}
		log.Printf("\n%s", string(dump))
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if api.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}
		log.Printf("\n%s", string(dump))
	}

	return resp, nil

}

// copyHeader copies all headers for `source` and sets them on `target`.
// based on https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}
