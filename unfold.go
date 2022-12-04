package main

import (
	"encoding/json"
	"fmt"
)

func Unfold(row string) ([]string, error) {
	b := []byte(row)
	var s []string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("unfold error: %w", err)
	}
	return s, nil
}
