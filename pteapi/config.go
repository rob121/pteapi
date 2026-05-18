package pteapi

import "strings"

const (
	DefaultBaseURL   = "https://app.plantoeat.com"
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

// Config holds credentials and HTTP options for the Plan to Eat client.
type Config struct {
	BaseURL   string
	Email     string
	Password  string
	UserAgent string
}

func (c Config) baseURL() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(c.BaseURL, "/")
}

func (c Config) userAgent() string {
	if strings.TrimSpace(c.UserAgent) == "" {
		return DefaultUserAgent
	}
	return c.UserAgent
}
