package fs

import (
	"fmt"
	"io"
	"os"
)

type FS struct{}

func (f *FS) Mkdir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	return nil
}

func (f *FS) Write(path string, contents io.Reader) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, contents); err != nil {
		return fmt.Errorf("failed to copy contents to file: %s", err)
	}

	return nil
}
