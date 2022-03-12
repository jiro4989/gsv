package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSV_ReadFold(t *testing.T) {
	tests := []struct {
		desc    string
		r       io.Reader
		lf      string
		want    []string
		wantErr bool
	}{
		{
			desc:    "正常系: カラムごとに取得できる",
			r:       strings.NewReader("a,b,c"),
			lf:      "\n",
			want:    []string{"a", "b", "c"},
			wantErr: false,
		},
		{
			desc:    "正常系: カンマが混在",
			r:       strings.NewReader(`"a,b","c"`),
			lf:      "\n",
			want:    []string{`a,b`, `c`},
			wantErr: false,
		},
		{
			desc:    "正常系: 改行文字",
			r:       strings.NewReader("\"a\nb\",c,\"d\ne\""),
			lf:      "\n",
			want:    []string{`a\nb`, "c", `d\ne`},
			wantErr: false,
		},
		{
			desc:    "正常系: CRLF改行文字",
			r:       strings.NewReader("\"a\r\nb\",c,\"d\r\ne\""),
			lf:      "\r\n",
			want:    []string{`a\nb`, "c", `d\ne`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			c := NewCSV(tt.r, tt.lf)
			got, err := c.ReadFold()
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

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

func TestFoldUnfold(t *testing.T) {
	tests := []struct {
		desc string
		s    string
	}{
		{
			desc: "正常系: 変化しない",
			s: `a\
b\a\b\c\n
\あ\い\う\え\お`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := unfold(fold(tt.s))
			assert.Equal(tt.s, got)
		})
	}
}
