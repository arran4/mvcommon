package mvcommon

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// CommonPrefixSplit finds the common prefix of strings, stopping at a stopWords if encountered, removing trim characters
// from the start and end, and ensuring that the match is minMatch in size minimum.
func CommonPrefixSplit(names []string, stopWords []string, trim string, minMatch int) string {
	if len(names) == 0 {
		return ""
	}

	type Match struct {
		pos int
		len int
	}

	type MatchSummary struct {
		Matches    []Match
		AveragePos int
	}

	matchLookup := make(map[string]map[int][]Match)
	trimMap := maps.Collect(func(yield func(K string, V struct{}) bool) {
		for i := 0; i < len(trim); i++ {
			if !yield(trim[i:i+1], struct{}{}) {
				return
			}
		}
	})

	var best *MatchSummary

	for ni, name := range names {
		for i := 0; i < len(name); i++ {
			me := Match{pos: i, len: 1}
			substr := name[i : i+1]
			if _, ok := trimMap[substr]; ok {
				continue
			}
			skip := 0
			for _, stopWord := range stopWords {
				if len(name) > i+len(stopWord) && name[i:i+len(stopWord)] == stopWord {
					skip = max(len(stopWord), skip)
				}
			}
			if skip > 0 {
				i += skip - 1
				continue
			}
			_, ok := matchLookup[substr]
			if !ok {
				matchLookup[substr] = map[int][]Match{ni: {me}}
				continue
			}
			matchLookup[substr][ni] = append(matchLookup[substr][ni], me)
		}
	}

	for len(matchLookup) > 0 {
		nextLookup := make(map[string]map[int][]Match)
		for matchStr, prefixMap := range matchLookup {
			if len(prefixMap) < len(names) {
				continue
			}
			trimChar := false
			if _, ok := trimMap[matchStr[len(matchStr)-1:]]; ok {
				trimChar = true
			}
			if len(matchStr) >= minMatch && !trimChar {
				sum := 0
				firstMatch := slices.Collect(func(yield func(Match) bool) {
					for nameI := range names {
						matchesForEach := prefixMap[nameI]
						firstMatch := matchesForEach[0]
						sum += firstMatch.pos
						if !yield(firstMatch) {
							return
						}
					}
				})
				averagePos := (sum * 100) / len(names)
				if best == nil || best.AveragePos >= averagePos {
					best = &MatchSummary{
						Matches:    firstMatch,
						AveragePos: averagePos,
					}
				}
			}
			for nameI, nameMatchesForPrefix := range prefixMap {
				name := names[nameI]
				for i := range nameMatchesForPrefix {
					matchesForPrefix := nameMatchesForPrefix[i]
					skip := 0
					for _, stopWord := range stopWords {
						if len(name) > matchesForPrefix.pos+matchesForPrefix.len-1+len(stopWord) && name[matchesForPrefix.pos+matchesForPrefix.len-1:matchesForPrefix.pos-1+matchesForPrefix.len+len(stopWord)] == stopWord {
							skip = max(len(stopWord), skip)
						}
					}
					if skip > 0 {
						continue
					}
					matchesForPrefix.len++
					if matchesForPrefix.pos+matchesForPrefix.len >= len(name) {
						continue
					}
					substr := name[matchesForPrefix.pos : matchesForPrefix.pos+matchesForPrefix.len]
					_, ok := nextLookup[substr]
					if !ok {
						nextLookup[substr] = map[int][]Match{nameI: {matchesForPrefix}}
						continue
					}
					nextLookup[substr][nameI] = append(nextLookup[substr][nameI], matchesForPrefix)
				}
			}
		}

		if len(nextLookup) == 0 {
			break
		}
		matchLookup = nextLookup
	}

	if best == nil {
		return ""
	}
	prefix := names[0][best.Matches[0].pos : best.Matches[0].pos+best.Matches[0].len]

	// Trim spaces and clean up the prefix
	return strings.Trim(prefix, trim)
}

// MoveFilesToFolder moves files into a specified folder. In dry-run mode, it only prints actions.
func MoveFilesToFolder(folder string, files []string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[Dry Run] Would create folder: %s\n", folder)
	} else {
		// Create folder if it doesn't exist
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("failed to create folder %s: %v", folder, err)
		}
	}

	// Move files into the folder
	for _, file := range files {
		destPath := filepath.Join(folder, filepath.Base(file))
		if dryRun {
			fmt.Printf("[Dry Run] Would move %s -> %s\n", file, destPath)
		} else {
			if err := os.Rename(file, destPath); err != nil {
				return fmt.Errorf("failed to move file %s: %v", file, err)
			}
			fmt.Printf("Moved %s -> %s\n", file, destPath)
		}
	}
	return nil
}
