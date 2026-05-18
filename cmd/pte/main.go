package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/rob121/pteapi/pteapi"
)

func main() {
	email := flag.String("email", envOr("PTE_EMAIL", ""), "Plan to Eat account email")
	password := flag.String("password", envOr("PTE_PASSWORD", ""), "Plan to Eat account password")
	userAgent := flag.String("user-agent", envOr("PTE_USER_AGENT", pteapi.DefaultUserAgent), "HTTP User-Agent header")
	baseURL := flag.String("base-url", envOr("PTE_BASE_URL", pteapi.DefaultBaseURL), "Plan to Eat app base URL")
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	cfg := pteapi.Config{
		BaseURL:   *baseURL,
		Email:     *email,
		Password:  *password,
		UserAgent: *userAgent,
	}

	client, err := pteapi.New(cfg)
	if err != nil {
		fatal(err)
	}

	if err := client.Login(); err != nil {
		fatal(fmt.Errorf("login: %w", err))
	}

	switch flag.Arg(0) {
	case "login":
		fmt.Println("logged in")
	case "recipes":
		recipes, err := client.ListRecipes()
		if err != nil {
			fatal(err)
		}
		emit(recipes)
	case "recipe":
		if flag.NArg() < 2 {
			fatal(fmt.Errorf("usage: pte recipe <id>"))
		}
		id, err := strconv.ParseInt(flag.Arg(1), 10, 64)
		if err != nil {
			fatal(fmt.Errorf("invalid recipe id: %w", err))
		}
		recipe, err := client.GetRecipe(id)
		if err != nil {
			fatal(err)
		}
		emit(recipe)
	case "plan", "planner":
		plan, err := client.MealPlan()
		if err != nil {
			fatal(err)
		}
		emit(plan)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `pte — Plan to Eat CLI

Usage:
  pte login
  pte recipes
  pte recipe <id>
  pte plan

Credentials (flags or environment):
  -email      PTE_EMAIL
  -password   PTE_PASSWORD
  -user-agent PTE_USER_AGENT
  -base-url   PTE_BASE_URL

`)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func emit(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
