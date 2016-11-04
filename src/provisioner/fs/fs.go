package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	FileModeRootReadWrite           = 0644
	FileModeRootReadWriteExecutable = 0744
)

type FS struct{}

func (fs *FS) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f *FS) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (f *FS) Mkdir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	return nil
}

func (f *FS) Write(path string, contents io.Reader, perm os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, contents); err != nil {
		return fmt.Errorf("failed to copy contents to file: %s", err)
	}

	return nil
}
