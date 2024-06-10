package filesystem

import (
	go_errors "errors"
	scanner_errors "ibm/container_cryptography_scanner/scanner/errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// A simple interface for a function to walk directories
type SimpleWalkDirFunc func(path string) error

// Filesystem interface is mainly used to interact with all types of possible data source (e.g. directories, docker images etc.); for images this represents a squashed layer
type Filesystem interface {
	WalkDir(fn SimpleWalkDirFunc) (err error)         // Walk the full filesystem using the SimpleWalkDirFunc fn
	ReadFile(path string) (content []byte, err error) // Read a specific file with a path from root of the filesystem
	GetConfig() (config v1.Config, ok bool)           // Get a config of this filesystem in container image format (if it exists)
}

// Simple plain filesystem that is constructed from the directory
type PlainFilesystem struct { // implements Filesystem
	rootPath string
}

func NewPlainFilesystem(rootPath string) PlainFilesystem {
	return PlainFilesystem{
		rootPath: rootPath,
	}
}

func (plainFilesystem PlainFilesystem) WalkDir(fn SimpleWalkDirFunc) error {
	return filepath.WalkDir(plainFilesystem.rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		err = fn(path)

		if go_errors.Is(err, scanner_errors.ErrParsingFailedAlthoughChecked) {
			slog.Warn(err.Error())
			return nil
		} else {
			return err
		}
	})
}

func (plainFilesystem PlainFilesystem) ReadFile(path string) ([]byte, error) {
	contentBytes, err := os.ReadFile(path)
	return contentBytes, err
}

// A plain directory does not have filesystem, so we return an empty object and false
func (plainFilesystem PlainFilesystem) GetConfig() (config v1.Config, ok bool) {
	return v1.Config{}, false
}
