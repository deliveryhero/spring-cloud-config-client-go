package configserverclient

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/deliveryhero/spring-cloud-config-client-go/src/logging"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var _ Client = (*client)(nil) // compile time proof

// Client defines api actions for Checkout API.
type Client interface {
	Get(context.Context, string, string) (map[string]any, error)
}

type client struct {
	url string
	log.Logger
	logger   logging.Logger
	client   *http.Client
	username string
	password string
}

// New creates a new client for Checkout API.
// defaultRetryWaitMin = 1 * time.Second.
// defaultRetryWaitMax = 30 * time.Second.
func New(options ...func(*client)) Client {
	c := &client{}

	for _, o := range options {
		o(c)
	}
	return c
}

// WithUsername sets the username for the client. Required.
func WithUsername(username string) func(*client) {
	return func(s *client) {
		s.username = username // pragma: allowlist username
	}
}

// WithPassword sets the password for the client. Required.
func WithPassword(password string) func(*client) {
	return func(s *client) {
		s.password = password // pragma: allowlist password
	}
}

// WithURL sets the url for the client. Required.
func WithURL(url string) func(*client) {
	return func(s *client) {
		s.url = url
	}
}

func withRetry(retry int, logger logging.Logger) func(*client) {
	return func(s *client) {
		rclient := retryablehttp.NewClient()
		rclient.RetryMax = retry
		rclient.Logger = logger

		s.logger = logger
		s.client = httptrace.WrapClient(rclient.StandardClient())
	}
}

// WithRetry3 sets the retry count 3 for the client and logger. Minimum 1 retry option required.
func WithRetry3(logger logging.Logger) func(*client) {
	return withRetry(3, logger)
}

// WithRetry5 sets the retry count 5 for the client and logger. Minimum 1 retry option required.
func WithRetry5(logger logging.Logger) func(*client) {
	return withRetry(5, logger)
}

// WithRetry10 sets the retry count 10 for the client and logger. Minimum 1 retry option required.
func WithRetry10(logger logging.Logger) func(*client) {
	return withRetry(10, logger)
}

// Get implements create payment checkout api client functionality.
func (c *client) Get(
	ctx context.Context,
	application string,
	environment string,
) (map[string]any, error) {
	url, err := url.Parse(c.url + "/" + application + "/" + environment)
	if err != nil {
		if c.logger != nil {
			c.logger.ErrorContext(ctx, "[ConfigServerClient].[Get]",
				"method", "ConfigServerClient.Get",
				"url-parse-error", err)
		}
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] url parse error")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		if c.logger != nil {
			c.logger.ErrorContext(ctx, "[ConfigServerClient].[Get]",
				"method", "CheckoutAPI.Get",
				"new-request-error", err)
		}
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] new request error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("connection", "keep-alive")

	res, err := c.client.Do(req)
	if err != nil {
		if c.logger != nil {
			c.logger.ErrorContext(ctx, "[ConfigServerClient].[Get]",
				"method", "ConfigServerClient.Get",
				"get-error", err)
		}
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] get error")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		if c.logger != nil {
			c.logger.ErrorContext(ctx, "[ConfigServerClient].[Get]",
				"method", "ConfigServerClient.Get",
				"read-body-error", err)
		}
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] read body error")
	}

	// 201 Created success response
	if res.StatusCode == http.StatusOK {
		body := map[string]any{}
		if err := json.Unmarshal(resBody, &body); err != nil {
			return nil, err
		}
		return body, nil
	}

	resBodyString := string(resBody)
	if c.logger != nil {
		c.logger.ErrorContext(ctx, "[ConfigServerClient].[Get]",
			"method", "ConfigServerClient.Get",
			"status-code", res.StatusCode,
			"status-message", res.Status,
			"error-response", resBodyString)
	}
	return nil, errors.New(
		"[ConfigServerClient].[Get] status:" + res.Status + " message: " + resBodyString,
	)
}
