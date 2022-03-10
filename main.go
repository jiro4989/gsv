package main

import (
	"encoding/csv"
	"io"
	"os"
)

type exitCode int

const (
	appName             = "gsv"
	exitCodeOK exitCode = iota
	exitCodeArgsErr
	exitCodeReadFoldErr
	exitCodeReadUnfoldErr
)

func main() {
	os.Exit(int(Main()))
}

func Main() exitCode {
	l := NewLogger(appName, os.Stdout, os.Stderr)
	args, err := ParseArgs()
	if err != nil {
		l.Err(err)
		return exitCodeArgsErr
	}

	if args.Ungsv {
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
