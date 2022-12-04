package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnfold2(t *testing.T) {
	tests := []struct {
		desc    string
		row     string
		want    []string
		wantErr bool
	}{
		{
			desc: "ok: json array to string array",
			row:  `["hello","world"]`,
			want: []string{"hello", "world"},
		},
		{
			desc: "ok: json array to string array",
			row:  `["hello\nworld"]`,
			want: []string{"hello\nworld"},
		},
		{
			desc: "ok: json array to string array",
			row:  `["hello\\nworld"]`,
			want: []string{"hello\\nworld"},
		},
		{
			desc:    "ng: fail",
			row:     `"`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := Unfold(tt.row)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			assert.Equal(tt.want, got)
			assert.NoError(err)
		})
	}
}
