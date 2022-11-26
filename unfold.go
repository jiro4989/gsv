package main

import "encoding/json"

func Unfold(row string) ([]string, error) {
	b := []byte(row)
	var s []string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}
