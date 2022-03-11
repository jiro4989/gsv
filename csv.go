package main

import (
	"encoding/csv"
	"io"
	"strings"
)

type CSV struct {
	reader *csv.Reader
	lf     string
}

func NewCSV(r io.Reader, lf string) *CSV {
	c := csv.NewReader(r)
	return &CSV{
		reader: c,
		lf:     lf,
	}
}

func (c *CSV) ReadFold() ([]string, error) {
	return c.read(fold)
}

func (c *CSV) ReadUnfold() ([]string, error) {
	return c.read(unfold)
}

func (c *CSV) read(fn func(string, string) string) ([]string, error) {
	row, err := c.reader.Read()
	if err != nil {
		return nil, err
	}
	result := make([]string, len(row))
	for i, cell := range row {
		result[i] = fn(cell, c.lf)
	}
	return result, nil
}

// fold folds a csv multiline-cell to oneline.
func fold(s, lf string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, lf, `\n`)
	return s
}

func unfold(s, lf string) string {
	r := strings.NewReader(s)
	var result []string
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}
		sch := string(ch)
		if sch == `\` {
			ch2, _, err := r.ReadRune()
			if err == io.EOF {
				result = append(result, sch)
				break
			}
			sch2 := string(ch2)
			if sch2 == `\` {
				// append '\' when ch + ch2 == '\\'
				result = append(result, `\`)
			} else if string(ch2) == "n" {
				// append '\n' when ch + ch2 == '\n'
				result = append(result, lf)
			} else {
				result = append(result, sch)
				result = append(result, sch2)
			}
			continue
		}
		result = append(result, sch)
	}
	return strings.Join(result, "")
}
