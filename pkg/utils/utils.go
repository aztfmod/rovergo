//
// Rover - utils and shared functions
// * Common functions and stuff that doesn't have a better home
//

package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func GenerateRandomGUID() string {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// ReadYamlFile finds extension of the given fileName
// Calculates fileName without extension
// Adds yaml and yml extensions to the fileName
// Gets the content of fileName.yaml or fileName.yml
// Returns the content of the file
func ReadYamlFile(filePath string) ([]byte, string, error) {
	extension := filepath.Ext(filePath)

	fileName := filepath.Base(filePath)

	if extension != "" && extension != ".yaml" && extension != ".yml" {
		return nil, fileName, fmt.Errorf("file extension must be .yaml or .yml")
	}

	filePathWithoutExtension := strings.TrimSuffix(filePath, extension)

	var err error
	var fileContent []byte

	if FileExists(filePathWithoutExtension + ".yaml") {
		fileName = filePathWithoutExtension + ".yaml"
	} else if FileExists(filePathWithoutExtension + ".yml") {
		fileName = filePathWithoutExtension + ".yml"
	} else {
		return nil, "", fmt.Errorf("could not find file '%s.yaml' or '%s.yml'", filePathWithoutExtension, filePathWithoutExtension)
	}

	fileContent, err = os.ReadFile(fileName)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file %s", fileName)
	}
	return fileContent, fileName, nil

}

var CurrentCustomCommandsAndGroupsYamlFilePath = ""

var SymphonyYamlFilePath = os.Getenv("HOME") + "/.rover/symphony.yaml"

func GetCustomCommandsAndGroupsYamlFilePath() string {
	return CurrentCustomCommandsAndGroupsYamlFilePath
}
