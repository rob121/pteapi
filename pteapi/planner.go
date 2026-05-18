package pteapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// MealPlan fetches the planner calendar page and parses planned meals.
func (c *Client) MealPlan() (*MealPlan, error) {
	if !c.LoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}

	html, err := c.fetchHTML("/planner")
	if err != nil {
		return nil, fmt.Errorf("fetch planner: %w", err)
	}

	return parsePlannerHTML(html)
}

func parsePlannerHTML(html []byte) (*MealPlan, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	plan := &MealPlan{}

	if data := doc.Find(".plannerData").First(); data.Length() > 0 {
		if raw, ok := data.Attr("data-recipes"); ok {
			_ = json.Unmarshal([]byte(raw), &plan.RecipeIDs)
		}
		if raw, ok := data.Attr("data-events"); ok {
			_ = json.Unmarshal([]byte(raw), &plan.EventIDs)
		}
	}

	doc.Find("#planner .planner-items").Each(func(_ int, item *goquery.Selection) {
		eventID, _ := strconv.ParseInt(item.AttrOr("data-id", ""), 10, 64)
		if eventID == 0 {
			return
		}
		recipeID, _ := strconv.ParseInt(item.AttrOr("data-recipe-id", ""), 10, 64)

		kind := "recipe"
		for _, class := range strings.Fields(item.AttrOr("class", "")) {
			if strings.HasPrefix(class, "kind-") {
				kind = strings.TrimPrefix(class, "kind-")
				break
			}
		}

		title := strings.TrimSpace(item.Find("a.title").First().Text())
		heading := strings.TrimSpace(item.Find(".item-heading").First().Text())
		servings := item.AttrOr("data-recipe-servings", item.Find("a.title").AttrOr("data-servings", ""))

		date := ""
		section := ""
		item.ParentsFiltered("td.date").Each(func(_ int, td *goquery.Selection) {
			date = td.AttrOr("data-date", "")
		})
		item.ParentsFiltered("div.time").Each(func(_ int, slot *goquery.Selection) {
			if section == "" {
				section = slot.AttrOr("data-section", "")
			}
		})

		plan.Items = append(plan.Items, PlanItem{
			EventID:  eventID,
			RecipeID: recipeID,
			Title:    title,
			Date:     date,
			Section:  section,
			Heading:  heading,
			Servings: servings,
			Kind:     kind,
		})
	})

	return plan, nil
}
