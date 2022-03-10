package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFold(t *testing.T) {
	tests := []struct {
		desc string
		s    string
		want string
	}{
		{
			desc: "正常系: 正常に折り畳める",
			s: `a
b
c`,
			want: `a\nb\nc`,
		},
		{
			desc: "正常系: 正常に折り畳める2",
			s: `a
あ\n
c\`,
			want: `a\nあ\\n\nc\\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := fold(tt.s)
			assert.Equal(tt.want, got)
		})
	}
}

func TestUnfold(t *testing.T) {
	tests := []struct {
		desc string
		s    string
		want string
	}{
		{
			desc: "正常系: 正常に展開できる",
			s:    `a\\b\\\\あ\nい\\nn\c`,
			want: `a\b\\あ
い\nn\c`,
		},
		{
			desc: "正常系: 変化なし",
			s:    `abcdeこんにちは`,
			want: `abcdeこんにちは`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := unfold(tt.s)
			assert.Equal(tt.want, got)
		})
	}
}
