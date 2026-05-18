//go:build integration

package pteapi

import (
	"os"
	"testing"
)

func TestIntegrationLoginAndFetch(t *testing.T) {
	email := os.Getenv("PTE_EMAIL")
	password := os.Getenv("PTE_PASSWORD")
	if email == "" || password == "" {
		t.Skip("set PTE_EMAIL and PTE_PASSWORD for integration tests")
	}

	client, err := New(Config{Email: email, Password: password})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Login(); err != nil {
		t.Fatalf("login: %v", err)
	}

	recipes, err := client.ListRecipes()
	if err != nil {
		t.Fatalf("list recipes: %v", err)
	}
	if len(recipes) == 0 {
		t.Fatal("expected at least one recipe")
	}

	recipe, err := client.GetRecipe(recipes[0].ID)
	if err != nil {
		t.Fatalf("get recipe: %v", err)
	}
	if recipe.Title == "" {
		t.Fatal("recipe title empty")
	}

	plan, err := client.MealPlan()
	if err != nil {
		t.Fatalf("meal plan: %v", err)
	}
	if plan == nil {
		t.Fatal("nil meal plan")
	}
}
