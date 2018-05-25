package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"os"
)

var (
	/**
	 * execution mode
	 */

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

	isTerminal = terminal.IsTerminal
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sqlfmt [parameter]... \n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	executionMode := parseMode()
	if err := exec(executionMode)(); err != nil {
		fmt.Println(err)
	}
}

type execFunc func() error

func parseMode() int {
	if !isTerminal(0) {
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

const (
	modeUnknown = 0
	modeDialog  = iota + 1
	modeCommand
	modeFile
	modePipe
)

func exec(mode int) execFunc {
	switch mode {
	case modeDialog:
		return dialogMode
	case modeCommand:
		return commandMode
	case modeFile:
		return fileMode
	case modePipe:
		return pipeMode
	default:
		return usageMode
	}
}

func dialogMode() error {
	f := func() string {
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
	for {
		str := f()
		if err := fmtSQL(str, os.Stdout, modeDialog); err != nil {
			return err
		}
	}
}

func commandMode() error {
	return fmtSQL(*sqlOpt, os.Stdout, modeCommand)
}

func fileMode() error {
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

func pipeMode() error {
	return fmtSQL(string(pipeStdIn), os.Stdout, modePipe)
}

func usageMode() error {
	usage()
	return nil
}

func file(f *os.File, w io.Writer) error {
	return fmtFile(f.Name(), f, w)
}
