package timer

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/cubny/httpqueue/internal/app/timer"
)

var (
	// ErrRetryableRequestFailure indicates that the request failed but most likely because of retryable errors such as http.StatusInternalServerError
	ErrRetryableRequestFailure = errors.New("HTTP request failed with a retryable error")
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically, so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)
	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically, so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by net/http when the
	// TLS certificate is not trusted. This error isn't typed
	// specifically, so we resort to matching on the error string.
	notTrustedErrorRe = regexp.MustCompile(`certificate is not trusted`)
)

// Client is a specialized HTTP Client for calling timer webhook.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	httpClient := http.DefaultClient
	return &Client{httpClient: httpClient}
}

// Shoot calls the timer's webhook and return ErrRetryableRequestFailure if the request fails because of a retryable
// reason, such as HTTP status code 500.
func (c *Client) Shoot(ctx context.Context, timer *timer.Timer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, timer.URL.String(), strings.NewReader(""))
	if err != nil {
		return err
	}

	resp, doErr := c.httpClient.Do(req)

	// Note: HTTP status code 429 response may include additional data in the header such as `Retry-After` which
	// should be respected in order to implement a "polite" client.
	// TODO: implement special treatment for 429 client error.
	shouldRetry, checkErr := isHTTPStatusCodeRetryable(resp, doErr)
	switch {
	case shouldRetry && checkErr != nil:
		return fmt.Errorf("http request failed: %v, %w", checkErr, ErrRetryableRequestFailure)
	case shouldRetry && checkErr == nil:
		return fmt.Errorf("http request failed: %v, %w", checkErr, ErrRetryableRequestFailure)
	case !shouldRetry && checkErr == nil:
		return nil
	}

	return nil
}

// isHTTPStatusCodeRetryable is aware of retry-ability of the request based on the response.
// The content of this function is inspired from HashiCorp's go-retryablehttp https://github.com/hashicorp/go-retryablehttp
func isHTTPStatusCodeRetryable(resp *http.Response, err error) (bool, error) {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if notTrustedErrorRe.MatchString(v.Error()) {
				return false, v
			}
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, v
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from client.
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side.
	// This will catch invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	return false, nil
}
