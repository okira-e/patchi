package utils

import (
	"encoding/json"
	"github.com/Okira-E/patchi/safego"
	"os"
)

// ReadJSONFile reads a JSON file and unmarshals it into a value reference.
func ReadJSONFile(filePath string, valRef any) safego.Option[error] {
	fileContent, err := os.ReadFile(filePath)

	err = json.Unmarshal(fileContent, valRef)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// WriteToJSONFile writes a value to a JSON file.
func WriteToJSONFile(filePath string, content any) safego.Option[error] {
	file, err := os.Create(filePath)
	defer file.Close()

	fileContent, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return safego.Some(err)
	}

	_, err = file.Write(fileContent)
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}
