package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "replyzer_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		filePath string
		data     []byte
		wantErr  bool
	}{
		{
			name:     "write to existing directory",
			filePath: filepath.Join(tempDir, "test.txt"),
			data:     []byte("test data"),
			wantErr:  false,
		},
		{
			name:     "write to nested directory",
			filePath: filepath.Join(tempDir, "nested", "dir", "test.txt"),
			data:     []byte("nested test data"),
			wantErr:  false,
		},
		{
			name:     "write empty data",
			filePath: filepath.Join(tempDir, "empty.txt"),
			data:     []byte(""),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteToFile(tt.filePath, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created and has correct content
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("File was not created: %s", tt.filePath)
					return
				}

				content, err := os.ReadFile(tt.filePath)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
					return
				}

				if string(content) != string(tt.data) {
					t.Errorf("File content mismatch. Expected: %s, Got: %s", string(tt.data), string(content))
				}
			}
		})
	}
}

func TestWriteToFilePermissions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replyzer_perm_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "perm_test.txt")
	data := []byte("permission test")

	err = WriteToFile(filePath, data)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	expectedPerm := os.FileMode(0644)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("File permissions mismatch. Expected: %v, Got: %v", expectedPerm, info.Mode().Perm())
	}
}