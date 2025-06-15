package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/arran4/mvcommon"
	"os"
	"strings"
)

func main() {
	// Define flags
	flagSet := flag.NewFlagSet("mvcommon", flag.ExitOnError)
	stopWords := mvcommon.DefaultStopWords
	defaultStopWords := true
	interactiveFlag := flagSet.Bool("interactive", false, "Enable interactive mode for file selection")
	dryRunFlag := flagSet.Bool("dry-run", false, "Perform a dry run without moving files")
	trimFlag := flagSet.String("trim", mvcommon.DefaultTrim, "Characters to trim")
	minMatchFlag := flagSet.Int("min", 3, "Minimum size of common segment")
	flagSet.Func("stopword", "Stopword to stop common prefix detection, defaults:[`"+strings.Join(stopWords, "`,`")+"`]", func(s string) error {
		if defaultStopWords {
			stopWords = []string{s}
			defaultStopWords = false
		} else {
			stopWords = append(stopWords, s)
		}
		return nil
	})
	flagSet.Usage = func() {
		Usage(stopWords, *trimFlag)
	}
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	// Remaining arguments are file names
	files := flagSet.Args()
	if len(files) < 2 {
		fmt.Println("Error: At least two files required")
		Usage(stopWords, *trimFlag)
		os.Exit(1)
	}

	var folderName string

	if *interactiveFlag {
		files, folderName = interactiveFileSelection(files, stopWords, *trimFlag, *minMatchFlag)
	} else {
		folderName = mvcommon.CommonPrefixSplit(files, stopWords, *trimFlag, *minMatchFlag)
	}
	if folderName == "" {
		fmt.Println("Error: No common prefix found! Exiting")
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files selected. Exiting.")
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

func interactiveFileSelection(files []string, stopWords []string, trim string, minMatch int) ([]string, string) {
	reader := bufio.NewReader(os.Stdin)
	selectedFiles := files
	for {

		// Print files with indices
		fmt.Println("\nInteractive Mode Enabled:")
		// Find common prefix
		folderName := mvcommon.CommonPrefixSplit(selectedFiles, stopWords, trim, minMatch)
		if folderName == "" {
			fmt.Println("Error: No common prefix found!")
		} else {
			fmt.Printf("Will move the files to %q\n\n", folderName)
		}

		fmt.Println("For the following files:")
		for i, file := range selectedFiles {
			fmt.Printf("%d. %s\n", i+1, file)
		}

		var nextSelectedFiles = make([]string, 0, len(selectedFiles))

		// Prompt user for confirmation or range input
		fmt.Println("\nEnter file numbers to include (e.g., 1,2,3 or 1-3,5-6) or press 'a' to confirm all, 'r' to reset:")
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

func Usage(stopWords []string, trimFlag string) {
	fmt.Println("Usage: mvcommon [-stopword=<stopword:`" + strings.Join(stopWords, "`,`") + "`>] [-trim=<trim:" + trimFlag + ">] [-min=3] [-dry-run] [-interactive] <file1> <file2> ...")
}
