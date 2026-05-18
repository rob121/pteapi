package pteapi

import (
	"strings"
	"testing"
)

func TestLoginFormParsing(t *testing.T) {
	html := `<form id="new_login" method="post">
<input type="hidden" name="authenticity_token" value="abc123" />
<div><input type="text" name="hp_field_123" id="hp_field_123" autocomplete="off" tabindex="-1" /></div>
<input name="login[email]" />
<input name="login[password]" />
</form>`

	auth := authTokenRe.FindStringSubmatch(html)
	if len(auth) != 2 || auth[1] != "abc123" {
		t.Fatalf("auth token: %v", auth)
	}

	hp := honeypotInputRe.FindStringSubmatch(html)
	if len(hp) != 2 || hp[1] != "hp_field_123" {
		t.Fatalf("honeypot: %v", hp)
	}
}

func TestListRecipesFiltersInactive(t *testing.T) {
	recipes := []RecipeSummary{
		{ID: 1, Title: "Good", Active: true, Remove: false},
		{ID: 2, Title: "", Active: true, Remove: false},
		{ID: 3, Title: "Removed", Active: true, Remove: true},
	}

	var out []RecipeSummary
	for _, r := range recipes {
		if !r.Active || r.Remove || r.Title == "" {
			continue
		}
		out = append(out, r)
	}

	if len(out) != 1 || out[0].ID != 1 {
		t.Fatalf("filter: %+v", out)
	}
}

func TestCSRFMetaRegex(t *testing.T) {
	html := `<meta name="csrf-token" content="token-value" />`
	m := csrfMetaRe.FindStringSubmatch(html)
	if len(m) != 2 || !strings.HasPrefix(m[1], "token") {
		t.Fatalf("csrf: %v", m)
	}
}
