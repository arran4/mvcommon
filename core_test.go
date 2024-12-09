package mvcommon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommonPrefixSplit(t *testing.T) {
	tests := []struct {
		name          string
		names         []string
		stopWords     []string
		trim          string
		minimumLength int
		expected      string
	}{
		{
			name:          "CommonPrefix_WithStopWord_AndTrimmedSpaces",
			names:         []string{"Report 234 - Draft1.txt", "Report 234 - Draft2.txt", "Report 234 - Final.txt"},
			stopWords:     []string{" - ", "] ", "["},
			trim:          "_- ",
			minimumLength: 3,
			expected:      "Report 234",
		},
		{
			name:          "CommonPrefix_WithStopWord_AndTrimmedSpaces Stop Words!",
			names:         []string{"Artist - Album I", "Artist - Album II"},
			stopWords:     []string{" - ", "] ", "["},
			trim:          "_- ",
			minimumLength: 3,
			expected:      "Artist",
		},
		{
			name:          "CommonPrefix_WithStopWords_AndBracketsInNames",
			names:         []string{"[Draft] Report 234.txt", "[For Review a] Report 234 - Version 2.txt", "[Final] Report 234.txt"},
			stopWords:     []string{" - ", "] ", "["},
			minimumLength: 3,
			trim:          "_- ",
			expected:      "Report 234",
		},
		{
			name:          "CommonPrefix_WithStopWords_MinimumLengthZero",
			names:         []string{"[Draft] Report 234.txt", "[For a Review] Report 234 - Version 2.txt", "[Final] Report 234.txt"},
			stopWords:     []string{" - ", "] ", "["},
			minimumLength: 0,
			trim:          "_- ",
			expected:      "a",
		},
		{
			name:          "CommonPrefix_NamesWithUnderscore_NoStopWord",
			names:         []string{"file_one.txt", "file_two.txt", "file_three.txt"},
			stopWords:     []string{" - ", "] ", "["},
			minimumLength: 3,
			trim:          "_- ",
			expected:      "file",
		},
		{
			name:          "CommonPrefix_NamesWithUnderscore_CommonPrefixApple",
			names:         []string{"apple_pie.txt", "apple_crumble.txt", "apple_sauce.txt"},
			stopWords:     []string{" - ", "] ", "["},
			trim:          "_- ",
			minimumLength: 3,
			expected:      "apple",
		},
		{
			name:      "CommonPrefix_NamesWithFinalStopWord",
			names:     []string{"data.csv", "data_final.csv", "data_backup.csv"},
			stopWords: []string{"final"},
			trim:      "_- ",
			expected:  "data",
		},
		{
			name:          "CommonPrefix_ExactMatchWithNoStopWord",
			names:         []string{"common_prefix.txt", "common_prefix_log.txt"},
			stopWords:     []string{" - ", "] ", "["},
			trim:          "_- ",
			minimumLength: 3,
			expected:      "common_prefix",
		},
		{
			name:          "CommonPrefix_EmptyNames_ReturnsEmptyString",
			names:         []string{},
			trim:          "_- ",
			minimumLength: 3,
			stopWords:     []string{" - ", "] ", "["},
			expected:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CommonPrefixSplit(test.names, test.stopWords, test.trim, test.minimumLength)
			if result != test.expected {
				t.Errorf("CommonPrefixSplit(%v, %#v, %#v, %d) got %q; want %q", test.names, test.stopWords, test.trim, test.minimumLength, result, test.expected)
			}
		})
	}
}

func TestDryRunMoveFiles(t *testing.T) {
	// Simulate file moving
	files := []string{"testfile1.txt", "testfile2.txt"}
	folder := "testfolder"

	err := MoveFilesToFolder(folder, files, true)
	if err != nil {
		t.Fatalf("Dry run failed: %v", err)
	}

	// No real changes should be made
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		t.Errorf("Dry run created folder %s unexpectedly", folder)
	}
}

func TestMoveFilesActual(t *testing.T) {
	t.Skip()
	// Setup temp files
	tempDir := t.TempDir()
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.txt"),
	}
	for _, file := range files {
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}

	// Move files
	folder := filepath.Join(tempDir, "output")
	err := MoveFilesToFolder(folder, files, false)
	if err != nil {
		t.Fatalf("MoveFilesToFolder failed: %v", err)
	}

	// Check folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		t.Errorf("Folder %s was not created", folder)
	}

	// Check files are moved
	for _, file := range files {
		destFile := filepath.Join(folder, filepath.Base(file))
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Errorf("File %s was not moved", destFile)
		}
	}
}
