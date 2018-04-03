package main

import (
	"fmt"
	"os"
	"bufio"
	"flag"
	"io"
	"bytes"
	"io/ioutil"
	"errors"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	/**
	 * execution mode
	 */
	executionMode = 0

	pipeStdIn []byte

	/**
	 * -sql: 1 command mode. output formatted sql to stdout. or use with -f, -ff
	 *  -d : interactive mode.
	 *  -f : file mode. rewrite selected '.go' file.
	 */
	sqlOpt         = flag.String("sql", "", "command line mode.")
	interactiveOpt = flag.Bool("i", false, "interactive mode.")
	fileOpt        = flag.String("f", "", "file mode.")
	writeOpt       = flag.Bool("w", false, "overwrite in the case of '-f', '-r'.")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sqlfmt [parameter]... \n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	executionMode = parseMode()
	if err := exec(); err != nil {
		fmt.Println(err)
	}
}

func parseMode() int {
	if !terminal.IsTerminal(0) {
		pipeStdIn, _ = ioutil.ReadAll(os.Stdin)
		return modePipe
	}

	if *fileOpt != "" {
		return modeFile
	}
	if *sqlOpt != "" {
		return modeCommand
	}
	if *interactiveOpt {
		return modeDialog
	}
	return modeUnknown
}

func exec() error {
	switch executionMode {
	case modeDialog:
		return dialogMode()
	case modeCommand:
		return fmtSQL(*sqlOpt, os.Stdout, modeCommand)
	case modeFile:
		if err := execFileMode(); err != nil {
			return errors.New(fmt.Sprintf("%s: %s", *fileOpt, err.Error()))
		}
		return nil
	case modePipe:
		return fmtSQL(string(pipeStdIn), os.Stdout, modePipe)
	case modeFiles:
	default:
		// helo mode
		usage()
	}
	return errors.New("Please input parameter.")
}

func dialogMode() error {
	for {
		str := dialog()
		if err := fmtSQL(str, os.Stdout, modeDialog); err != nil {
			return err
		}
	}
}

const (
	modeUnknown = 0
	modeDialog  = iota + 1
	modeCommand
	modeFile
	modePipe
	modeFiles
)

func dialog() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("('3') type sql to here: ")
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if text := scanner.Text(); text != "" {
			return text
		}
	}
	return ""
}

func execFileMode() error {
	f, err := os.Open(*fileOpt)
	if err != nil {
		return err
	}
	defer f.Close()
	buff := new(bytes.Buffer)
	if err := file(f, buff); err != nil {
		return err
	}
	if *writeOpt {
		if err := ioutil.WriteFile(*fileOpt, buff.Bytes(), os.ModePerm); err != nil {
			return err
		}
	} else {
		if _, err := os.Stdout.WriteString(buff.String()); err != nil {
			return err
		}
	}
	return nil
}

func file(f *os.File, w io.Writer) error {
	return fmtFile(f.Name(), f, w)
}
