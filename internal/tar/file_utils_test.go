package tar

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTaritAndUntar(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tarit_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file structure
	testDir := filepath.Join(tempDir, "foo-archive")
	err = os.Mkdir(testDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testFile := filepath.Join(testDir, "test_file.txt")
	err = os.WriteFile(testFile, []byte("Test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test Tarit
	tarFile := filepath.Join(tempDir, "test.tar.gz")
	err = Tarit(testDir, tarFile)
	if err != nil {
		t.Fatalf("Tarit failed: %v", err)
	}

	// Check if the tar file was created
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		t.Fatalf("Tar file was not created")
	}

	// Test Untar
	extractDir := filepath.Join(tempDir, "bar-archive-extracted")
	err = os.Mkdir(extractDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create extraction directory: %v", err)
	}

	tarFileHandle, err := os.Open(tarFile)
	if err != nil {
		t.Fatalf("Failed to open tar file: %v", err)
	}
	defer tarFileHandle.Close()

	err = Untar(context.Background(), tarFileHandle, extractDir)
	if err != nil {
		t.Fatalf("Untar failed: %v", err)
	}

	// Check if the extracted file exists and has the correct content
	extractedFile := filepath.Join(extractDir, "test_file.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "Test content" {
		t.Errorf("Extracted file content mismatch. Expected 'Test content', got '%s'", string(content))
	}
}
