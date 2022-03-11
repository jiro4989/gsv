package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFold(t *testing.T) {
	tests := []struct {
		desc string
		s    string
		lf   string
		want string
	}{
		{
			desc: "正常系: 正常に折り畳める",
			s: `a
b
c`,
			lf:   "\n",
			want: `a\nb\nc`,
		},
		{
			desc: "正常系: 正常に折り畳める2",
			s: `a
あ\n
c\`,
			lf:   "\n",
			want: `a\nあ\\n\nc\\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := fold(tt.s, tt.lf)
			assert.Equal(tt.want, got)
		})
	}
}

func TestUnfold(t *testing.T) {
	tests := []struct {
		desc string
		s    string
		lf   string
		want string
	}{
		{
			desc: "正常系: 正常に展開できる",
			s:    `a\\b\\\\あ\nい\\nn\c`,
			lf:   "\n",
			want: `a\b\\あ
い\nn\c`,
		},
		{
			desc: "正常系: 変化なし",
			s:    `abcdeこんにちは`,
			lf:   "\n",
			want: `abcdeこんにちは`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := unfold(tt.s, tt.lf)
			assert.Equal(tt.want, got)
		})
	}
}
