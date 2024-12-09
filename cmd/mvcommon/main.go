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

	// Find common prefix
	folderName := mvcommon.CommonPrefixSplit(files, stopWords, *trimFlag, *minMatchFlag)
	if folderName == "" {
		fmt.Println("No common prefix found. Exiting.")
		os.Exit(1)
	}

	if *interactiveFlag {
		files = interactiveFileSelection(files)
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

func interactiveFileSelection(files []string) []string {
	reader := bufio.NewReader(os.Stdin)

	// Print files with indices
	fmt.Println("\nInteractive Mode Enabled:")
	fmt.Println("The following files are detected:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	// Prompt user for confirmation or range input
	fmt.Println("\nEnter file numbers to include (e.g., 1,2,3 or 1-3,5-6) or press 'a' to confirm all:")
	for {
		fmt.Print("Your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "a" {
			return files // Confirm all files
		}

		selectedIndices, err := mvcommon.ParseNumberRanges(input, len(files))
		if err != nil {
			fmt.Println("Invalid input:", err)
			continue
		}

		var selectedFiles []string
		for _, idx := range selectedIndices {
			selectedFiles = append(selectedFiles, files[idx])
		}

		if len(selectedFiles) > 0 {
			fmt.Println("Selected files:")
			for _, file := range selectedFiles {
				fmt.Println(" -", file)
			}
			return selectedFiles
		}

		fmt.Println("No valid files selected. Try again.")
	}
}

func Usage(stopWords []string, trimFlag string) {
	fmt.Println("Usage: mvcommon [-stopword=<stopword:`" + strings.Join(stopWords, "`,`") + "`>] [-trim=<trim:" + trimFlag + ">] [-min=3] [-dry-run] [-interactive] <file1> <file2> ...")
}
