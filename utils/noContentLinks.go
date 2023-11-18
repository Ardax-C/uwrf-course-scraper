package utils

import (
	"encoding/json"
	"os"
)

func SaveNoContentLinksToFile(data []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	return err
}
