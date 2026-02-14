package main

import (
	"bufio"
	"fmt"
	"github.com/arran4/mvcommon"
	"os"
	"strings"
)

// Run is a subcommand `mvcommon`
//
// Flags:
//
//	stopWords:	--stopword		Stop word to stop common prefix detection
//	trim:		--trim			Characters to trim
//	minMatch:	--min			Minimum size of common segment
//	dryRun:		--dry-run		Perform a dry run without moving files
//	interactive:	--interactive	Enable interactive mode for file selection
//	files:		...				Files to move
func Run(stopWords string, trim string, minMatch int, dryRun bool, interactive bool, files ...string) {
	if len(files) < 2 {
		fmt.Println("Error: At least two files required")
		Usage()
		os.Exit(1)
	}

	var stopWordsSlice []string
	if stopWords != "" {
		stopWordsSlice = strings.Split(stopWords, ",")
	} else {
		stopWordsSlice = mvcommon.DefaultStopWords
	}

	var folderName string

	if interactive {
		files, folderName = interactiveFileSelection(files, stopWordsSlice, trim, minMatch)
	} else {
		folderName = mvcommon.CommonPrefixSplit(files, stopWordsSlice, trim, minMatch)
	}
	if folderName == "" {
		fmt.Println("Error: No common prefix found! Exiting")
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files selected. Exiting.")
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Creating folder: %s\n", folderName)
	} else {
		fmt.Printf("Creating folder: %s\n", folderName)
	}

	// Move files into the folder
	if err := mvcommon.MoveFilesToFolder(folderName, files, dryRun); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully.")
}

func interactiveFileSelection(files []string, stopWords []string, trim string, minMatch int) ([]string, string) {
	reader := bufio.NewReader(os.Stdin)
	selectedFiles := files
	for {

		// Print files with indices
		fmt.Println()
		fmt.Println("Interactive Mode Enabled:")
		// Find common prefix
		folderName := mvcommon.CommonPrefixSplit(selectedFiles, stopWords, trim, minMatch)
		if folderName == "" {
			fmt.Fprintln(os.Stderr, "Error: No common prefix found!")
		} else {
			fmt.Printf("Will move the files to %q\n", folderName)
			fmt.Println()
		}

		fmt.Println("For the following files:")
		for i, file := range selectedFiles {
			fmt.Printf("%d. %s\n", i+1, file)
		}

		var nextSelectedFiles = make([]string, 0, len(selectedFiles))

		// Prompt user for confirmation or range input
		fmt.Println()
		fmt.Println("Enter file numbers to include (e.g., 1,2,3 or 1-3,5-6) or press 'a' to confirm all, 'r' to reset:")
		for {
			fmt.Print("Your choice: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				panic(err)
			}
			input = strings.TrimSpace(input)

			if input == "a" {
				return selectedFiles, folderName // Confirm all files
			}

			if input == "r" {
				nextSelectedFiles = files
				break
			}

			selectedIndices, err := mvcommon.ParseNumberRanges(input, len(selectedFiles))
			if err != nil {
				fmt.Println("Invalid input:", err)
				continue
			}

			for _, idx := range selectedIndices {
				nextSelectedFiles = append(nextSelectedFiles, selectedFiles[idx])
			}

			if len(nextSelectedFiles) > 0 {
				break
			}

			fmt.Println("No valid files selected. Try again.")
		}
		selectedFiles = nextSelectedFiles
	}
}

func Usage() {
	stopWords := mvcommon.DefaultStopWords
	trimFlag := mvcommon.DefaultTrim
	fmt.Println("Usage: mvcommon [-stopword=<stopword:`" + strings.Join(stopWords, "`,`") + "`>] [-trim=<trim:" + trimFlag + ">] [-min=3] [-dry-run] [-interactive] <file1> <file2> ...")
}
