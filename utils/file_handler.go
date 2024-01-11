package utils

import (
	"encoding/json"
	"io/ioutil"
)

// SaveToFile takes a string (JSON data) and a filename, and writes the data to the file.
func SaveToFile(data, filename string) error {
	// Unmarshal the data to ensure it's valid JSON
	var jsonData interface{}
	err := json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		// Data is not valid JSON
		return err
	}

	// Re-marshal the data to format it as JSON
	formattedData, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	// Write the formatted data to the file
	err = ioutil.WriteFile(filename, formattedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
