package myutils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// function to create a hash of a string using sha256
func HashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// function Struct2json to convert a struct to json string
func Struct2json(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
