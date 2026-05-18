package pteapi

import (
	"fmt"
	"strconv"
)

// ListRecipes returns recipe summaries from the authenticated recipe book.
// Plan to Eat exposes this data at /api/v1/recipes/lite_plus for logged-in users.
func (c *Client) ListRecipes() ([]RecipeSummary, error) {
	if !c.LoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}

	var recipes []RecipeSummary
	if err := c.doJSON("GET", "/api/v1/recipes/lite_plus", &recipes); err != nil {
		return nil, err
	}

	out := make([]RecipeSummary, 0, len(recipes))
	for _, r := range recipes {
		if !r.Active || r.Remove || r.Title == "" {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

// GetRecipe returns full detail for a recipe ID.
func (c *Client) GetRecipe(id int64) (*Recipe, error) {
	if !c.LoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}
	if id <= 0 {
		return nil, fmt.Errorf("invalid recipe id %d", id)
	}

	var recipe Recipe
	path := "/api/v1/recipes/" + strconv.FormatInt(id, 10)
	if err := c.doJSON("GET", path, &recipe); err != nil {
		return nil, err
	}
	return &recipe, nil
}
