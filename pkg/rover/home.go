package rover

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed home/*
var homeDefaultDir embed.FS

const roverHomePath = ".rover"

var homeDir string

func initializeHomeDir() error {
	home, err := os.UserHomeDir()

	if err != nil {
		return errors.New("Unable to access user home directory")
	}

	roverhome := filepath.Join(home, roverHomePath)

	// Create directory when it not exists and populate with default files
	_, err = os.Stat(roverhome)
	if os.IsNotExist(err) {
		direrr := os.MkdirAll(roverhome, 0777)
		if direrr != nil {
			return fmt.Errorf("Failed to create %s directory", roverhome)
		}

		err := createDefaultContents(roverhome)
		if err != nil {
			return err
		}
	}
	homeDir = roverhome
	return nil
}

func HomeDirectory() (string, error) {
	if len(strings.TrimSpace(homeDir)) == 0 {
		err := initializeHomeDir()
		if err != nil {
			return "", err
		}
	}
	return homeDir, nil
}

func SetHomeDirectory(dir string) {
	homeDir = dir
}

func createDefaultContents(roverHomePath string) error {
	return fs.WalkDir(homeDefaultDir, "home", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Strip the top folder out of the path
		outPath := strings.Replace(path, "home/", "", -1)
		// Skip these files
		if outPath == "home" || outPath == "readme.md" {
			return nil
		}
		outPath = filepath.Join(roverHomePath, outPath)

		if entry.IsDir() {
			// create directories
			err = os.MkdirAll(outPath, 0777)

			if err != nil {
				return err
			}
		} else {
			// copy files from the embedded FS to the real filesystem
			bytes, err := homeDefaultDir.ReadFile(path)
			if err != nil {
				return err
			}

			err = os.WriteFile(outPath, bytes, 0777)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
