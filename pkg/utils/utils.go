//
// Rover - utils and shared functions
// * Common functions and stuff that doesn't have a better home
// * Ben C, May 2021
//

package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/aztfmod/rover/pkg/console"
)

// StringSliceDel deletes a specific index from a slide of strings
// Taken from https://yourbasic.org/golang/delete-element-slice/
func StringSliceDel(a []string, i int) []string {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
	a[len(a)-1] = ""     // Erase last element (write zero value).
	a = a[:len(a)-1]     // Truncate slice.
	return a
}

// CopyFile is a very simple file copy helper
func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	bytesWritten, err := io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	console.Debugf("Completed copying file '%s' to '%s' (%d bytes)", src, dest, bytesWritten)
	return nil
}

func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", errors.New("Unable to access user home directory")
	}

	roverhome := filepath.Join(home, "/.rover")

	direrr := os.MkdirAll(roverhome, 0777)

	if direrr != nil {
		return "", errors.New("Failed to create $home/.rover directory")
	}

	return roverhome, nil
}
