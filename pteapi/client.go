package pteapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var (
	csrfMetaRe      = regexp.MustCompile(`<meta name="csrf-token" content="([^"]+)"`)
	authTokenRe     = regexp.MustCompile(`name="authenticity_token" value="([^"]+)"`)
	honeypotInputRe = regexp.MustCompile(`<input type="text" name="([^"]+)" id="[^"]+" autocomplete="off" tabindex="-1"`)
)

// Client scrapes Plan to Eat using colly with a shared cookie jar.
type Client struct {
	cfg Config

	collector *colly.Collector

	mu         sync.Mutex
	loggedIn   bool
	csrfToken  string
	lastStatus int
}

// New creates a client. Call Login before other methods.
func New(cfg Config) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookie jar: %w", err)
	}

	c := colly.NewCollector(
		colly.UserAgent(cfg.userAgent()),
		colly.AllowedDomains("app.plantoeat.com", "plantoeat.com"),
	)

	c.SetCookieJar(jar)

	client := &Client{
		cfg:       cfg,
		collector: c,
	}

	c.OnResponse(func(r *colly.Response) {
		client.mu.Lock()
		client.lastStatus = r.StatusCode
		client.mu.Unlock()
		if token := csrfMetaRe.FindStringSubmatch(string(r.Body)); len(token) == 2 {
			client.mu.Lock()
			client.csrfToken = token[1]
			client.mu.Unlock()
		}
	})

	return client, nil
}

// LoggedIn reports whether Login completed successfully.
func (c *Client) LoggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.loggedIn
}

func (c *Client) baseURL() string {
	return c.cfg.baseURL()
}

func (c *Client) setLoggedIn(ok bool) {
	c.mu.Lock()
	c.loggedIn = ok
	c.mu.Unlock()
}

func (c *Client) getCSRFToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.csrfToken
}

func (c *Client) fetchHTML(path string) ([]byte, error) {
	return c.scrape(http.MethodGet, path, nil, "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
}

func (c *Client) doJSON(method, path string, dest any) error {
	body, err := c.scrape(method, path, nil, "application/json")
	if err != nil {
		return err
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func (c *Client) postForm(path string, form url.Values) ([]byte, error) {
	return c.scrape(http.MethodPost, path, strings.NewReader(form.Encode()), "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Origin":       c.baseURL(),
		"Referer":      mustJoin(c.baseURL(), "/login"),
	})
}

func (c *Client) scrape(method, path string, body io.Reader, accept string, extraHeaders ...map[string]string) ([]byte, error) {
	target := mustJoin(c.baseURL(), path)

	var (
		respBody []byte
		status   int
		scrapeErr error
	)

	clone := c.collector.Clone()
	clone.OnResponse(func(r *colly.Response) {
		status = r.StatusCode
		respBody = append([]byte(nil), r.Body...)
	})

	clone.OnError(func(r *colly.Response, err error) {
		scrapeErr = err
	})

	clone.OnRequest(func(r *colly.Request) {
		r.Method = method
		r.Headers.Set("Accept", accept)
		if token := c.getCSRFToken(); token != "" {
			r.Headers.Set("X-CSRF-Token", token)
		}
		if method == http.MethodPost || method == http.MethodPut || strings.Contains(accept, "application/json") {
			r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		}
		for _, hdrs := range extraHeaders {
			for k, v := range hdrs {
				r.Headers.Set(k, v)
			}
		}
	})

	switch method {
	case http.MethodGet:
		if err := clone.Visit(target); err != nil {
			return nil, err
		}
	case http.MethodPost:
		if body == nil {
			return nil, fmt.Errorf("POST %s requires body", path)
		}
		data, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		if err := clone.PostRaw(target, data); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported method %s", method)
	}

	clone.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}
	if status >= 400 {
		return nil, fmt.Errorf("%s %s: status %d: %s", method, path, status, strings.TrimSpace(string(respBody)))
	}
	return respBody, nil
}

func mustJoin(base, path string) string {
	u, err := url.JoinPath(base, path)
	if err != nil {
		panic(err)
	}
	return u
}
