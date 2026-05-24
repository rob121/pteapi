package pteapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rating is a recipe star rating. The API may return null, a number, or a string.
type Rating struct {
	Value *float64
}

// Float returns the rating value, or nil when unset.
func (r Rating) Float() *float64 {
	return r.Value
}

func (r *Rating) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		r.Value = nil
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			r.Value = nil
			return nil
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("rating string %q: %w", s, err)
		}
		r.Value = &f
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	r.Value = &f
	return nil
}

func (r Rating) MarshalJSON() ([]byte, error) {
	if r.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*r.Value)
}
