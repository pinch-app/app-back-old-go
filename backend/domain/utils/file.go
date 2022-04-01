package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// File  Model
type File struct {
	Key         string
	ContentType string
	Format      string
	Location    string
}

// Get file content
func (f *File) Content() (string, error) {
	file, err := ioutil.ReadFile(f.Path())
	if err != nil {
		return "", err
	}
	return string(file), nil
}

// open file
func (f *File) Open() (*os.File, error) {
	return os.Open(f.Path())
}

// Get local file path
func (f *File) Path() string {
	location := f.Location
	if location == "" {
		location = os.TempDir()
	}
	return filepath.Join(location, f.Name())
}

// Remove local file
func (f *File) Close() error {
	os.Remove(f.Path())
	return nil
}

// Get File Name with extension
func (f *File) Name() string {
	return f.Key + "." + f.Format
}
