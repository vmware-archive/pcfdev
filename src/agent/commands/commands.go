package commands

import (
	"io"
	"os"
)

//go:generate mockgen -package mocks -destination mocks/cmd_runner.go agent/commands CmdRunner
type CmdRunner interface {
	Run(command string, args ...string) error
	Output(command string, args ...string) (output []byte, err error)
}

//go:generate mockgen -package mocks -destination mocks/fs.go agent/commands FS
type FS interface {
	Mkdir(directory string) error
	Write(path string, contents io.Reader, perm os.FileMode) error
	Read(path string) (contents []byte, err error)
	Exists(path string) (bool, error)
}

