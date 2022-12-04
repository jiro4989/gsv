package main

import (
	"encoding/json"
	"fmt"
)

func Fold(row []string) (string, error) {
	b, err := json.Marshal(row)
	if err != nil {
		return "", fmt.Errorf("fold error: %w", err)
	}
	s := string(b)
	return s, nil
}
