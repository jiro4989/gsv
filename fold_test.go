package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFold2(t *testing.T) {
	tests := []struct {
		desc    string
		row     []string
		want    string
		wantErr bool
	}{
		{
			desc: "ok: simple row to json array",
			row:  []string{"hello", "world"},
			want: `["hello","world"]`,
		},
		{
			desc: "ok: multi-line row to json array",
			row:  []string{"hello\nworld"},
			want: `["hello\nworld"]`,
		},
		{
			desc: "ok: \\n row to json array",
			row:  []string{"hello\\nworld"},
			want: `["hello\\nworld"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := Fold(tt.row)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			assert.Equal(tt.want, got)
			assert.NoError(err)
		})
	}
}
