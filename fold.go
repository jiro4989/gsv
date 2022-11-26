package main

import "encoding/json"

func Fold(row []string) (string, error) {
	b, err := json.Marshal(row)
	if err != nil {
		return "", err
	}
	s := string(b)
	return s, nil
}
