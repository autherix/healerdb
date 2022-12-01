package myutils

import (
	"os"
)

// Function to get the current working directory, return the path as string and error
func GetCwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

// Function to get the path of the file, return the path parts(directories, filename, and extension seperately) as string and error
func GetPathParts(path string) (string, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return "", "", "", err
	}
	return fi.Name(), fi.Name(), fi.Name(), nil
}

// Functions to check whether a slice contains a string or not
func ContainsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
