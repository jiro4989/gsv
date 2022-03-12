package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type exitCode int

type Param struct {
	Ungsv  bool
	LF     string
	Output string
	Args   []string
}

type App struct {
	param  Param
	logger Logger
}

const (
	appName             = "gsv"
	version             = "dev"
	exitCodeOK exitCode = iota
	exitCodeArgsErr
	exitCodeReadFoldErr
	exitCodeReadUnfoldErr
	exitCodeOpenFileErr
)

var (
	param     Param
	convertLF = map[string]string{
		"lf":   "\n",
		"crlf": "\r\n",
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&param.Ungsv, "ungsv", "u", false, "unfold csv rows")
	rootCmd.Flags().StringVarP(&param.LF, "linefeed", "l", "lf", "input text line feed character. [lf | crlf]")
	rootCmd.Flags().StringVarP(&param.Output, "output", "o", "", "output file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "'gsv' transforms a multi-line CSV into one-line CSV to make it easier to 'grep'",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		param.Args = args
		Main(param)
	},
}

func Main(p Param) exitCode {
	a := NewApp(p)

	if err := p.Validate(); err != nil {
		a.logger.Err(err)
		return exitCodeArgsErr
	}

	// open file when args have a file
	var inputFile *os.File
	defer inputFile.Close()
	if 0 < len(p.Args) {
		var err error
		inputFile, err = os.Open(p.Args[0])
		if err != nil {
			a.logger.Err(err)
			return exitCodeOpenFileErr
		}
	} else {
		inputFile = os.Stdin
	}

	// create a file when an output file path is empty
	var outputFile *os.File
	defer outputFile.Close()
	if 0 < len(p.Output) {
		var err error
		outputFile, err = os.Create(p.Output)
		if err != nil {
			a.logger.Err(err)
			return exitCodeOpenFileErr
		}
	} else {
		outputFile = os.Stdout
	}

	if p.Ungsv {
		if err := a.readUnfoldAndWrite(inputFile, outputFile); err != nil {
			a.logger.Err(err)
			return exitCodeReadUnfoldErr
		}
		return exitCodeOK
	}

	if err := a.readFoldAndWrite(inputFile, outputFile); err != nil {
		a.logger.Err(err)
		return exitCodeReadFoldErr
	}
	return exitCodeOK
}

func NewApp(p Param) *App {
	l := NewLogger(appName, os.Stdout, os.Stderr)
	return &App{
		param:  p,
		logger: l,
	}
}

func (a *App) readFoldAndWrite(r io.Reader, w io.Writer) error {
	fn := func(c *CSV) ([]string, error) {
		return c.ReadFold()
	}
	return a.readAndWrite(r, w, fn)
}

func (a *App) readUnfoldAndWrite(r io.Reader, w io.Writer) error {
	fn := func(c *CSV) ([]string, error) {
		return c.ReadUnfold()
	}
	return a.readAndWrite(r, w, fn)
}

func (a *App) readAndWrite(r io.Reader, w io.Writer, fn func(c *CSV) ([]string, error)) error {
	c := NewCSV(r, convertLF[a.param.LF])
	useCRLF := a.param.LF == "crlf"
	for {
		row, err := fn(c)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		w := csv.NewWriter(w)
		w.UseCRLF = useCRLF
		if err := w.Write(row); err != nil {
			return err
		}
		w.Flush()
	}
	return nil
}

func (p *Param) Validate() error {
	_, ok := convertLF[p.LF]
	if !ok {
		return fmt.Errorf("'%s' of '--linefeed' is not supported", p.LF)
	}
	return nil
}
