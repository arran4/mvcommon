package skill

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsLocalPath(source string) bool {
	return strings.HasPrefix(source, ".") || strings.HasPrefix(source, "/") || strings.HasPrefix(source, "~")
}

// FetchEmbedded copies the embedded official skill directory to the destination
func FetchEmbedded(source string, dest string) error {
	// Source is either 'embedded' or 'official'
	// Extract from embed.FS

	// Read embedded file
	data, err := OfficialSkills.ReadFile("skills/mvcommon/SKILL.md")
	if err != nil {
		return fmt.Errorf("failed to read embedded official skill: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write file
	destFile := filepath.Join(dest, "SKILL.md")
	if err := os.WriteFile(destFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write embedded skill to %s: %w", destFile, err)
	}

	return nil
}

// FetchLocal copies a local skill directory to the destination
func FetchLocal(source string, dest string) error {
	// First, verify that source exists and is a directory
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("could not stat local source: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("local source is not a directory: %s", source)
	}

	// Create destination directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Simple recursive copy
	return copyDir(source, dest)
}

func copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// preserve basic permissions
	info, err := os.Stat(src)
	if err == nil {
		os.Chmod(dst, info.Mode())
	}

	return nil
}

// FetchRemote fetches a remote skill (e.g. GitHub repository)
// source format: owner/repo or github.com/owner/repo
func FetchRemote(source string, dest string) (string, error) {
	// Simple resolution for GitHub repositories
	// If it doesn't contain a slash, or is malformed, we error
	parts := strings.Split(source, "/")
	var owner, repo string

	if len(parts) == 2 {
		owner, repo = parts[0], parts[1]
	} else if len(parts) >= 3 && parts[0] == "github.com" {
		owner, repo = parts[1], parts[2]
	} else {
		return "", fmt.Errorf("unsupported remote source format (use owner/repo): %s", source)
	}

	// Fetch main branch tarball
	// E.g., https://api.github.com/repos/owner/repo/tarball/main
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tarball", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Read github token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download skill: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download skill, HTTP %d", resp.StatusCode)
	}

	// Find the resolved commit SHA from the ETag or headers if available
	// GitHub API usually returns the SHA in the ETag header like W/"sha"
	commitSHA := ""
	etag := resp.Header.Get("ETag")
	if etag != "" {
		commitSHA = strings.Trim(etag, "W/\"")
	}

	// Decompress gzip
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read gzip stream: %w", err)
	}
	defer gzr.Close()

	// Extract tar
	tr := tar.NewReader(gzr)

	// Create dest directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar header: %w", err)
		}

		// Security: Check for path traversal (e.g., ../)
		if strings.Contains(header.Name, "..") {
			return "", fmt.Errorf("invalid path in tar archive: %s", header.Name)
		}

		// The first directory in the tarball is usually owner-repo-sha/
		// We want to strip the first component
		parts := strings.SplitN(header.Name, "/", 2)
		if len(parts) < 2 {
			continue // Skip the root directory entry
		}

		targetPath := filepath.Join(dest, parts[1])

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return "", err
			}
			out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return "", err
			}
			out.Close()
		}
	}

	return commitSHA, nil
}
