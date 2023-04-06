package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDotDir(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal("cannot get working directory: ", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}

	// Test when .ggit folder exists
	if err := os.Mkdir(".ggit", 0755); err != nil {
		t.Errorf("Error creating .ggit folder: %v", err)
	}

	dotDir, err := getDotDir()
	if err != nil {
		t.Errorf("error in get dor folder")
	}
	if dotDir != filepath.Join(tempDir, ".ggit") {
		t.Errorf("Expected .ggit, got %s", dotDir)
	}

	// change back to original folder in order to succeed in deleting temp folder
	if err := os.Chdir(originalDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}
}

func TestGetDotDirFailed(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal("cannot get working directory: ", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}

	_, _ = getDotDir()

	// change back to original folder in order to succeed in deleting temp folder
	if err := os.Chdir(originalDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}
}

// TODO
func TestAdd(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal("cannot get working directory: ", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}

	// Test when .ggit folder exists
	if err := os.Mkdir(".ggit", 0755); err != nil {
		t.Errorf("Error creating .ggit folder: %v", err)
	}

	// os.NewFile()

	// dotDir, err := Add("file1")
	if err != nil {
		t.Errorf("Error finding .ggit folder: %v", err)
	}

	// change back to original folder in order to succeed in deleting temp folder
	if err := os.Chdir(originalDir); err != nil {
		t.Fatal("cannot change folder: ", err)
	}
}
