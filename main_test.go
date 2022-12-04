package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testOutputDir = "testdata/out"
)

func TestMain(m *testing.M) {
	testBefore()
	exitCode := m.Run()
	testAfter()
	os.Exit(exitCode)
}

func testBefore() {
	if err := os.Mkdir(testOutputDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func testAfter() {
	if err := os.RemoveAll(testOutputDir); err != nil {
		panic(err)
	}
}

func TestMain2(t *testing.T) {
	tests := []struct {
		desc        string
		p           Param
		want        exitCode
		wantStdout  string
		wantOutFile string
	}{
		{
			desc: "ok: multi line csv",
			p: Param{
				Ungsv:    false,
				LF:       "lf",
				Output:   "",
				Encoding: "utf8",
				Args:     []string{"testdata/sample1.csv"},
			},
			want: exitCodeOK,
			wantStdout: `["Language","Word","Note"]
["English","Hello\nWorld","note"]
["Japanese","こんにちは\nこんばんは","メモ"]
["English","John\nRose",""]
["Japanese","太郎\n花子",""]
`,
		},
		{
			desc: "ok: LF is used when Ungsv flag was true",
			p: Param{
				Ungsv:    false,
				LF:       "crlf",
				Output:   "",
				Encoding: "utf8",
				Args:     []string{"testdata/sample1.csv"},
			},
			want: exitCodeOK,
			wantStdout: `["Language","Word","Note"]
["English","Hello\nWorld","note"]
["Japanese","こんにちは\nこんばんは","メモ"]
["English","John\nRose",""]
["Japanese","太郎\n花子",""]
`,
		},
		{
			desc: "ok: write output to file",
			p: Param{
				Ungsv:    false,
				LF:       "lf",
				Output:   testOutputDir + "/1.csv",
				Encoding: "utf8",
				Args:     []string{"testdata/sample1.csv"},
			},
			want:        exitCodeOK,
			wantStdout:  "",
			wantOutFile: testOutputDir + "/1.csv",
		},
		{
			desc: "ok: use Ungsv",
			p: Param{
				Ungsv:    true,
				LF:       "lf",
				Output:   "",
				Encoding: "utf8",
				Args:     []string{"testdata/sample3_gsved_utf8_unix.txt"},
			},
			want: exitCodeOK,
			wantStdout: `Language,Word,Note
English,"Hello
World",note
Japanese,"こんにちは
こんばんは",メモ
English,"John
Rose",
Japanese,"太郎
花子",
`,
		},
		{
			desc: "ok: reading sjis csv",
			p: Param{
				Ungsv:    false,
				LF:       "crlf",
				Output:   "",
				Encoding: "sjis",
				Args:     []string{"testdata/sample4_sjis_crlf.txt"},
			},
			want: exitCodeOK,
			wantStdout: `["a","b"]
["日本","Japan"]
["あ","い\nう"]
`,
		},
		{
			desc: "err: LF is invalid",
			p: Param{
				Ungsv:    false,
				LF:       "sushi",
				Output:   "",
				Encoding: "utf8",
				Args:     []string{"testdata/sample1.csv"},
			},
			want: exitCodeArgsErr,
		},
		{
			desc: "err: unsupported encoding",
			p: Param{
				Ungsv:    false,
				LF:       "lf",
				Output:   "",
				Encoding: "sushi",
				Args:     []string{"testdata/sample1.csv"},
			},
			want: exitCodeArgsErr,
		},
		{
			desc: "err: input file is not found",
			p: Param{
				Ungsv:    false,
				LF:       "lf",
				Output:   "",
				Encoding: "utf8",
				Args:     []string{"sushi.txt"},
			},
			want: exitCodeOpenFileErr,
		},
		{
			desc: "err: output directory is not found",
			p: Param{
				Ungsv:    false,
				LF:       "lf",
				Output:   "/a/b/c/d.csv",
				Encoding: "utf8",
				Args:     []string{"testdata/sample1.csv"},
			},
			want: exitCodeOpenFileErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			// capture stdout
			r, w, err := os.Pipe()
			assert.NoError(err)

			stdout := os.Stdout
			os.Stdout = w

			got := Main(tt.p)
			os.Stdout = stdout
			w.Close()

			assert.Equal(tt.want, got)
			var output bytes.Buffer
			io.Copy(&output, r)
			assert.Equal(tt.wantStdout, output.String())

			if tt.wantOutFile != "" {
				_, err := os.Stat(tt.wantOutFile)
				exists := !os.IsNotExist(err)
				assert.True(exists)
			}
		})
	}
}
