package pteapi

import (
	"os"
	"testing"
)

func TestParsePlannerHTML(t *testing.T) {
	html, err := os.ReadFile("../.tmp/planner.html")
	if err != nil {
		t.Skip("planner fixture missing; run integration probe first")
	}

	plan, err := parsePlannerHTML(html)
	if err != nil {
		t.Fatal(err)
	}

	if len(plan.EventIDs) == 0 {
		t.Fatal("expected event ids")
	}
	if len(plan.Items) == 0 {
		t.Fatal("expected planner items")
	}

	item := plan.Items[0]
	if item.EventID != 76529158 {
		t.Fatalf("event id: got %d want 76529158", item.EventID)
	}
	if item.RecipeID != 49145900 {
		t.Fatalf("recipe id: got %d want 49145900", item.RecipeID)
	}
	if item.Title != "Cottage Cheese Flatbread" {
		t.Fatalf("title: got %q", item.Title)
	}
	if item.Date != "2026-05-19" {
		t.Fatalf("date: got %q want 2026-05-19", item.Date)
	}
	if item.Section != "dinner" {
		t.Fatalf("section: got %q want dinner", item.Section)
	}
}
