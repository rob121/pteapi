package pteapi

import (
	"fmt"
	"net/url"
	"strings"
)

// Login authenticates with email and password, retaining session cookies.
func (c *Client) Login() error {
	if c.cfg.Email == "" || c.cfg.Password == "" {
		return fmt.Errorf("email and password are required")
	}

	loginHTML, err := c.fetchHTML("/login")
	if err != nil {
		return fmt.Errorf("fetch login page: %w", err)
	}

	authMatch := authTokenRe.FindStringSubmatch(string(loginHTML))
	if len(authMatch) != 2 {
		return fmt.Errorf("authenticity_token not found on login page")
	}

	form := url.Values{}
	form.Set("authenticity_token", authMatch[1])
	form.Set("login[email]", c.cfg.Email)
	form.Set("login[password]", c.cfg.Password)
	if m := honeypotInputRe.FindStringSubmatch(string(loginHTML)); len(m) == 2 {
		form.Set(m[1], "")
	}

	respBody, err := c.postForm("/login", form)
	if err != nil {
		return fmt.Errorf("post login: %w", err)
	}

	// Successful login redirects away from /login (often to /recipes).
	if strings.Contains(string(respBody), `id="pg-login"`) {
		return fmt.Errorf("login failed: invalid credentials or blocked request")
	}

	c.setLoggedIn(true)
	return nil
}
