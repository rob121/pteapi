# Plan to Eat Api

Go client for [Plan to Eat](https://www.plantoeat.com). There is no public API; this package logs in via the web app, keeps session cookies, and fetches recipes and meal plans using [colly](https://github.com/gocolly/colly).

## Install

From source:

```bash
git clone https://github.com/rob121/pteapi.git
cd pteapi
go build -o pte ./cmd/pte
```

Or install the CLI:

```bash
go install github.com/rob121/pteapi/cmd/pte@latest
```

## Credentials

Set environment variables or pass flags on every command:

| Flag | Environment | Default |
|------|-------------|---------|
| `-email` | `PTE_EMAIL` | *(required)* |
| `-password` | `PTE_PASSWORD` | *(required)* |
| `-user-agent` | `PTE_USER_AGENT` | Chrome on macOS |
| `-base-url` | `PTE_BASE_URL` | `https://app.plantoeat.com` |

Optional: copy `pte.yaml.example` to `pte.yaml` (gitignored) as a local reference.

```bash
export PTE_EMAIL='you@example.com'
export PTE_PASSWORD='your-password'
```

## CLI

All data commands print indented JSON to stdout.

```bash
pte login              # verify credentials; prints "logged in"
pte recipes            # list recipe summaries
pte recipe <id>        # full recipe by ID
pte plan               # meal planner (alias: pte planner)
```

### `pte login`

```text
logged in
```

### `pte recipes`

Returns an array of `RecipeSummary` (inactive, removed, and untitled recipes are omitted).

```json
[
  {
    "id": 49145900,
    "title": "Cottage Cheese Flatbread",
    "course_title": "Main Course",
    "total_time": 35,
    "photo_url": "https://plantoeat.s3.amazonaws.com/recipes/49145900/...-large.jpg",
    "rating": null,
    "owned": true,
    "queued": false,
    "active": true,
    "remove": false,
    "stats_total_events": 1,
    "updated_at": "2026-05-18T14:44:50-04:00",
    "created_at": "2026-05-18T14:44:50-04:00"
  },
  {
    "id": 49146133,
    "title": "Strawberry Cottage Cheese Salad",
    "course_title": "Salad",
    "total_time": 5,
    "photo_url": "https://example.com/.../salad.jpg",
    "rating": null,
    "owned": true,
    "queued": true,
    "active": true,
    "remove": false,
    "stats_total_events": 0,
    "updated_at": "2026-05-18T14:59:14-04:00",
    "created_at": "2026-05-18T14:59:11-04:00"
  }
]
```

### `pte recipe <id>`

Returns a single `Recipe` with ingredients and nutrition fields.

```json
{
  "id": 49145900,
  "title": "Cottage Cheese Flatbread",
  "course_title": "Main Course",
  "total_time": 35,
  "photo_url": "https://plantoeat.s3.amazonaws.com/recipes/49145900/...-large.jpg",
  "rating": null,
  "owned": true,
  "queued": false,
  "active": true,
  "remove": false,
  "stats_total_events": 1,
  "updated_at": "2026-05-18T14:44:50-04:00",
  "created_at": "2026-05-18T14:44:50-04:00",
  "description": "Make this delicious and high protein Cottage Cheese Flatbread...",
  "source": "https://theproteinchef.co/cottage-cheese-flatbread-recipe/",
  "url": "https://theproteinchef.co/cottage-cheese-flatbread-recipe/",
  "servings": 1,
  "yield": "1 Wrap",
  "prep_time": 5,
  "cook_time": 30,
  "ingredient_titles": [
    "Cottage Cheese",
    "Garlic Powder",
    "Italian Seasoning",
    "Oregano",
    "Whole Eggs"
  ],
  "ingredients": [
    {
      "id": 595319417,
      "title": "Whole Eggs",
      "amount": "2",
      "unit": "Large",
      "note": "",
      "position": 0,
      "amount_float": 2,
      "metric_amount": 2,
      "metric_unit": "Large"
    }
  ],
  "private": false,
  "draft": false
}
```

### `pte plan`

Returns a `MealPlan` parsed from the planner calendar HTML for the currently visible date range.

```json
{
  "recipe_ids": [49145900],
  "event_ids": [76529158],
  "items": [
    {
      "event_id": 76529158,
      "recipe_id": 49145900,
      "title": "Cottage Cheese Flatbread",
      "date": "2026-05-19",
      "section": "dinner",
      "heading": "Leftover",
      "servings": "1",
      "kind": "recipe"
    }
  ]
}
```

| Field | Meaning |
|-------|---------|
| `event_id` | Planner event ID |
| `recipe_id` | Recipe ID (when `kind` is `recipe`) |
| `date` | `YYYY-MM-DD` |
| `section` | `breakfast`, `lunch`, `dinner`, `xtra`, or `notes` |
| `heading` | Label such as `Leftover` |
| `kind` | `recipe`, `ingredient`, `note`, etc. |

## Library

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rob121/pteapi/pteapi"
)

func main() {
	client, err := pteapi.New(pteapi.Config{
		Email:    os.Getenv("PTE_EMAIL"),
		Password: os.Getenv("PTE_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}
	if err := client.Login(); err != nil {
		panic(err)
	}

	// []RecipeSummary
	recipes, err := client.ListRecipes()
	if err != nil {
		panic(err)
	}

	// *Recipe — full detail for one ID
	recipe, err := client.GetRecipe(recipes[0].ID)
	if err != nil {
		panic(err)
	}

	// *MealPlan
	plan, err := client.MealPlan()
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(recipe)
	_ = enc.Encode(plan)
}
```

### API

| Method | Returns | Description |
|--------|---------|-------------|
| `New(cfg Config)` | `*Client` | Create client with cookie jar and colly collector |
| `(*Client) Login()` | `error` | POST login form; retains session cookies |
| `(*Client) LoggedIn()` | `bool` | Whether `Login` succeeded |
| `(*Client) ListRecipes()` | `[]RecipeSummary` | All active recipes in your book |
| `(*Client) GetRecipe(id int64)` | `*Recipe` | Single recipe by ID |
| `(*Client) MealPlan()` | `*MealPlan` | Planned meals on the visible calendar |

## Tests

```bash
go test ./...

# hits Plan to Eat with PTE_EMAIL / PTE_PASSWORD
go test -tags=integration ./pteapi -run Integration -v
```
