package main

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type exitCode int

type Param struct {
	Ungsv bool
}

const (
	appName             = "gsv"
	version             = "dev"
	exitCodeOK exitCode = iota
	exitCodeArgsErr
	exitCodeReadFoldErr
	exitCodeReadUnfoldErr
)

var (
	param Param
)

func init() {
	rootCmd.Flags().BoolVarP(&param.Ungsv, "ungsv", "u", false, "unfold csv rows")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "gsv",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Main(param)
	},
}

func Main(p Param) exitCode {
	l := NewLogger(appName, os.Stdout, os.Stderr)

	if p.Ungsv {
		if err := readUnfoldAndWrite(os.Stdin, os.Stdout); err != nil {
			l.Err(err)
			return exitCodeReadUnfoldErr
		}
		return exitCodeOK
	}

	if err := readFoldAndWrite(os.Stdin, os.Stdout); err != nil {
		l.Err(err)
		return exitCodeReadFoldErr
	}
	return exitCodeOK
}

func readFoldAndWrite(r io.Reader, w io.Writer) error {
	fn := func(c *CSV) ([]string, error) {
		return c.ReadFold()
	}
	return readAndWrite(r, w, fn)
}

func readUnfoldAndWrite(r io.Reader, w io.Writer) error {
	fn := func(c *CSV) ([]string, error) {
		return c.ReadUnfold()
	}
	return readAndWrite(r, w, fn)
}

func readAndWrite(r io.Reader, w io.Writer, fn func(c *CSV) ([]string, error)) error {
	c := NewCSV(r)
	for {
		row, err := fn(c)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		w := csv.NewWriter(w)
		if err := w.Write(row); err != nil {
			return err
		}
		w.Flush()
	}
	return nil
}
