package configserverclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var _ Client = (*client)(nil) // compile time proof

// Client defines api actions for Checkout API.
type Client interface {
	Get(context.Context, string, string, string) (map[string]any, error)
}

type client struct {
	url      string
	client   *http.Client
	username string
	password string
}

type ErrorResponse struct {
	error
	StatusCode int
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

func WithClient(httpClient *http.Client) func(*client) {
	return func(s *client) {
		s.client = httpClient
	}
}

// Get implements create payment checkout api client functionality.
func (c *client) Get(
	ctx context.Context,
	application string,
	environment string,
	label string,
) (map[string]any, error) {
	if c.client == nil {
		c.client = http.DefaultClient
	}

	url, err := url.Parse(c.url + "/" + application + "/" + environment + "/" + label)
	if err != nil {
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] url parse error")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] new request error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("connection", "keep-alive")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "[ConfigServerClient].[Get] get error")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
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
	return nil, ErrorResponse{
		error: errors.New(
			"[ConfigServerClient].[Get] status:" + res.Status + " message: " + resBodyString),
		StatusCode: res.StatusCode,
	}
}
