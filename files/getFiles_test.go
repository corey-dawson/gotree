package files

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// setupTestDir creates a temporary directory structure for testing.
// It returns the path to the temporary root and a cleanup function.
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper() // <-- Added: Marks this function as a test helper
	// 1. Create a temporary directory unique to this test
	tempDir, err := os.MkdirTemp("", "test-dir-tree")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Define the structure:
	// root/
	// ├── file1.txt
	// ├── sub1/
	// │   ├── file2.md
	// │   ├── sub2_empty/  <-- This should trigger the empty dir logic
	// │   └── .hidden_file <-- This is a hidden file and should be collected!
	// └── empty_dir/       <-- This should trigger the empty dir logic

	structure := []string{
		"file1.txt",
		"sub1/file2.md",
		"sub1/sub2_empty",
		"sub1/.hidden_file",
		"empty_dir",
	}

	for _, path := range structure {
		fullPath := filepath.Join(tempDir, path)

		// Check if it should be a directory or a file
		isDir := (path == "sub1/sub2_empty" || path == "empty_dir")

		if isDir {
			// Create the directory
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", fullPath, err)
			}
		} else {
			// It's a file, ensure the parent directory exists first
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				t.Fatalf("Failed to create parent dir %s: %v", parentDir, err)
			}
			// Create the file
			if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", fullPath, err)
			}
		}
	}

	// 2. Return the cleanup function to delete the temporary directory later
	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

// TestGetStructure runs the getStructure function against the temporary directory.
func TestGetStructure(t *testing.T) {

	t.Run("ShowAll", func(*testing.T) {
		tempDir, cleanup := setupTestDir(t)
		// IMPORTANT: Defer the cleanup so it runs when the test exits
		defer cleanup()

		// Expected results (paths must be absolute, matching getStructure output)
		sep := string(filepath.Separator)

		expectedFiles := []string{
			// Files
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "sub1", "file2.md"),
			filepath.Join(tempDir, "sub1", ".hidden_file"), // Note: The provided function includes hidden files

			// Empty Directories (must end with separator)
			filepath.Join(tempDir, "sub1", "sub2_empty") + sep,
			filepath.Join(tempDir, "empty_dir") + sep,
		}

		// 1. Run the function
		actualFiles := GetStructure(tempDir, true, -1, 1)

		// 2. Sort both slices for reliable comparison
		sort.Strings(expectedFiles)
		sort.Strings(actualFiles)

		// 3. Check length first
		if len(actualFiles) != len(expectedFiles) {
			t.Errorf("Mismatch in number of files found.\nExpected %d files:\n%v\nActual %d files:\n%v",
				len(expectedFiles), expectedFiles, len(actualFiles), actualFiles)
			return
		}

		// 4. Compare contents element by element
		for i := range actualFiles {
			if actualFiles[i] != expectedFiles[i] {
				t.Errorf("File at index %d is incorrect.\nExpected: %s\nActual: %s",
					i, expectedFiles[i], actualFiles[i])
			}
		}
	})

	t.Run("HideHidden", func(*testing.T) {
		tempDir, cleanup := setupTestDir(t)
		// IMPORTANT: Defer the cleanup so it runs when the test exits
		defer cleanup()

		// Expected results (paths must be absolute, matching getStructure output)
		// Note: .hidden_file should NOT appear in results when showAll is false
		sep := string(filepath.Separator)

		expectedFiles := []string{
			// Files (hidden files should be excluded)
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "sub1", "file2.md"),

			// Empty Directories (must end with separator)
			filepath.Join(tempDir, "sub1", "sub2_empty") + sep,
			filepath.Join(tempDir, "empty_dir") + sep,
		}

		// 1. Run the function
		actualFiles := GetStructure(tempDir, false, -1, 1)

		// 2. Verify hidden file is NOT in results
		hiddenFilePath := filepath.Join(tempDir, "sub1", ".hidden_file")
		for _, file := range actualFiles {
			if file == hiddenFilePath {
				t.Errorf("Hidden file .hidden_file should not appear when showAll is false, but found: %s", file)
			}
		}

		// 3. Sort both slices for reliable comparison
		sort.Strings(expectedFiles)
		sort.Strings(actualFiles)

		// 4. Check length first
		if len(actualFiles) != len(expectedFiles) {
			t.Errorf("Mismatch in number of files found.\nExpected %d files:\n%v\nActual %d files:\n%v",
				len(expectedFiles), expectedFiles, len(actualFiles), actualFiles)
			return
		}

		// 5. Compare contents element by element
		for i := range actualFiles {
			if actualFiles[i] != expectedFiles[i] {
				t.Errorf("File at index %d is incorrect.\nExpected: %s\nActual: %s",
					i, expectedFiles[i], actualFiles[i])
			}
		}
	})

	t.Run("MaxDepth", func(*testing.T) {
		tempDir, cleanup := setupTestDir(t)
		// IMPORTANT: Defer the cleanup so it runs when the test exits
		defer cleanup()

		// Test with maxDepth=1: should only show root level files and directories
		// but NOT recurse into subdirectories
		sep := string(filepath.Separator)

		expectedFiles := []string{
			// Files at root level (level 1)
			filepath.Join(tempDir, "file1.txt"),
			// Directories at level 1 should be included but not recursed
			filepath.Join(tempDir, "sub1") + sep,
			filepath.Join(tempDir, "empty_dir") + sep,
		}

		// 1. Run the function with maxDepth=1
		actualFiles := GetStructure(tempDir, true, 1, 1)

		// 2. Verify that files from deeper levels are NOT included
		deepFile := filepath.Join(tempDir, "sub1", "file2.md")
		deepDir := filepath.Join(tempDir, "sub1", "sub2_empty") + sep
		for _, file := range actualFiles {
			if file == deepFile {
				t.Errorf("File from level 2 should not appear when maxDepth=1, but found: %s", file)
			}
			if file == deepDir {
				t.Errorf("Directory from level 2 should not appear when maxDepth=1, but found: %s", file)
			}
		}

		// 3. Sort both slices for reliable comparison
		sort.Strings(expectedFiles)
		sort.Strings(actualFiles)

		// 4. Check length first
		if len(actualFiles) != len(expectedFiles) {
			t.Errorf("Mismatch in number of files found.\nExpected %d files:\n%v\nActual %d files:\n%v",
				len(expectedFiles), expectedFiles, len(actualFiles), actualFiles)
			return
		}

		// 5. Compare contents element by element
		for i := range actualFiles {
			if actualFiles[i] != expectedFiles[i] {
				t.Errorf("File at index %d is incorrect.\nExpected: %s\nActual: %s",
					i, expectedFiles[i], actualFiles[i])
			}
		}
	})
}
