package main

import (
	"flag"
	"fmt"
	"github.com/arran4/mvcommon"
	"os"
	"strings"
)

func main() {
	// Define flags
	flagSet := flag.NewFlagSet("mvcommon", flag.ExitOnError)
	stopWords := []string{
		" - ",
		"] ",
		"[",
	}
	defaultStopWords := true
	flagSet.Func("stopword", "Stopword to stop common prefix detection, defaults:[`"+strings.Join(stopWords, "`,`")+"`]", func(s string) error {
		if defaultStopWords {
			stopWords = []string{s}
			defaultStopWords = false
		} else {
			stopWords = append(stopWords, s)
		}
		return nil
	})
	dryRunFlag := flagSet.Bool("dry-run", false, "Perform a dry run without moving files")
	trimFlag := flagSet.String("trim", "-_ ", "Characters to trim")
	minMatchFlag := flagSet.Int("minimumMatch", 3, "Minimum size of common segment")
	flag.Parse()

	// Remaining arguments are file names
	files := flag.Args()
	if len(files) < 2 {
		fmt.Println("Usage: mvcommon [-stopword=<stopword: - >] [-trim=<trim:- _>] [-dry-run] <file1> <file2> ...")
		os.Exit(1)
	}

	// Find common prefix
	commonPrefix := mvcommon.CommonPrefixSplit(files, stopWords, *trimFlag, *minMatchFlag)
	if commonPrefix == "" {
		fmt.Println("No common prefix found. Exiting.")
		os.Exit(1)
	}

	// Clean the common prefix further for folder naming
	folderName := strings.ReplaceAll(commonPrefix, " ", "")
	if folderName == "" {
		fmt.Println("Invalid folder name after cleaning. Exiting.")
		os.Exit(1)
	}

	if *dryRunFlag {
		fmt.Printf("[Dry Run] Creating folder: %s\n", folderName)
	} else {
		fmt.Printf("Creating folder: %s\n", folderName)
	}

	// Move files into the folder
	if err := mvcommon.MoveFilesToFolder(folderName, files, *dryRunFlag); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully.")
}