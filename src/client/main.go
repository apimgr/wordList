package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apimgr/wordList/src/words"
)

var Version = "1.0.0"

func main() {
	name := filepath.Base(os.Args[0])
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		fmt.Fprintf(os.Stdout, "%s v%s\n\nUsage: %s passphrase [--words 6] [--separator -] [--lowercase]\n", name, Version, name)
		return
	}
	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%s v%s\n", name, Version)
		return
	}
	if os.Args[1] != "passphrase" {
		fmt.Fprintf(os.Stderr, "unknown command %q; use --help\n", os.Args[1])
		os.Exit(2)
	}
	flags := flag.NewFlagSet("passphrase", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	count := flags.Int("words", 6, "Number of words (1-20)")
	separator := flags.String("separator", "-", "Word separator (up to 8 characters)")
	lowercase := flags.Bool("lowercase", false, "Do not capitalize words")
	if err := flags.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}
	if *count < 1 || *count > 20 || len(*separator) > 8 {
		fmt.Fprintln(os.Stderr, "words must be 1-20 and separator at most 8 characters")
		os.Exit(2)
	}
	if err := words.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "word list unavailable")
		os.Exit(1)
	}
	passphrase, err := words.PassphraseSecure(*count, *separator, !*lowercase)
	if err != nil {
		fmt.Fprintln(os.Stderr, "secure random generation unavailable")
		os.Exit(1)
	}
	fmt.Println(passphrase)
}

