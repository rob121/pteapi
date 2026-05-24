package pteapi

import (
	"encoding/json"
	"testing"
)

func TestRatingUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want *float64
	}{
		{"null", `null`, nil},
		{"number", `4.5`, ptr(4.5)},
		{"string", `"3"`, ptr(3)},
		{"empty string", `""`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Rating
			if err := json.Unmarshal([]byte(tt.raw), &r); err != nil {
				t.Fatal(err)
			}
			if !floatPtrEqual(r.Value, tt.want) {
				t.Fatalf("got %v want %v", r.Value, tt.want)
			}
		})
	}
}

func TestRatingMarshal(t *testing.T) {
	var summary struct {
		Rating Rating `json:"rating"`
	}
	summary.Rating = Rating{Value: ptr(5)}
	b, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"rating":5}` {
		t.Fatalf("got %s", b)
	}
}

func ptr(f float64) *float64 { return &f }

func floatPtrEqual(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
