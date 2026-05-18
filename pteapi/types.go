package pteapi

import "time"

// RecipeSummary is a lightweight recipe entry from the recipe book.
type RecipeSummary struct {
	ID                  int64     `json:"id"`
	Title               string    `json:"title"`
	CourseTitle         string    `json:"course_title"`
	TotalTime           float64   `json:"total_time"`
	PhotoURL            string    `json:"photo_url"`
	Rating              *float64  `json:"rating"`
	Owned               bool      `json:"owned"`
	Queued              bool      `json:"queued"`
	Active              bool      `json:"active"`
	Remove              bool      `json:"remove"`
	StatsTotalEvents    int       `json:"stats_total_events"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedAt           time.Time `json:"created_at"`
}

// Ingredient is one line item on a recipe.
type Ingredient struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Amount       string  `json:"amount"`
	Unit         string  `json:"unit"`
	Note         string  `json:"note"`
	Position     int     `json:"position"`
	AmountFloat  float64 `json:"amount_float"`
	MetricAmount float64 `json:"metric_amount"`
	MetricUnit   string  `json:"metric_unit"`
}

// Recipe is full recipe detail from the authenticated API.
type Recipe struct {
	RecipeSummary
	Description      string       `json:"description"`
	Source           string       `json:"source"`
	URL              string       `json:"url"`
	Servings         int          `json:"servings"`
	Yield            string       `json:"yield"`
	PrepTime         float64      `json:"prep_time"`
	CookTime         float64      `json:"cook_time"`
	IngredientTitles []string     `json:"ingredient_titles"`
	Ingredients      []Ingredient `json:"ingredients"`
	Private          bool         `json:"private"`
	Draft            bool         `json:"draft"`
}

// PlanItem is a single meal on the planner calendar.
type PlanItem struct {
	EventID   int64  `json:"event_id"`
	RecipeID  int64  `json:"recipe_id,omitempty"`
	Title     string `json:"title"`
	Date      string `json:"date"`      // YYYY-MM-DD
	Section   string `json:"section"`   // breakfast, lunch, dinner, xtra, notes
	Heading   string `json:"heading"`   // e.g. Leftover
	Servings  string `json:"servings"`
	Kind      string `json:"kind"`      // recipe, ingredient, note, etc.
}

// MealPlan is the visible planner calendar state.
type MealPlan struct {
	RecipeIDs []int64    `json:"recipe_ids"`
	EventIDs  []int64    `json:"event_ids"`
	Items     []PlanItem `json:"items"`
}
