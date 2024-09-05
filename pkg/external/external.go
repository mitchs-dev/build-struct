package external

import (
	"errors"
	"os"

	"github.com/mitchs-dev/build-struct/pkg/builder"
)

/*
Call is the API-variant of the main function

This function requires structName. However, you can set filePath or yamlData to use one or the other.
*/
func Call(structName, filePath string, yamlData []byte) (string, error) {

	// Ensure that the structName is provided
	if structName == "" {
		return "", errors.New("structName is required")
	} else {
		// Set the struct name
		builder.StructName = structName
	}

	// Ensure that either filePath or fileData is provided
	if filePath == "" && len(yamlData) == 0 {
		return "", errors.New("filePath or fileData is required")
	}

	if filePath != "" && len(yamlData) != 0 {
		return "", errors.New("filePath and fileData cannot both be provided")
	}

	// If the file path is provided, read the file data
	if filePath != "" {
		// Make sure that the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", errors.New("error: The provided file does not exist")
		}

		// Make sure that the file is not empty
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", errors.New("error: The file is empty")
		}

		// Make sure that it's a file and not a directory
		fileInfo, err := os.Stat(filePath)
		if err != nil && fileInfo.IsDir() {
			return "", errors.New("error: The provided path is a directory")
		}

		// Read the file
		readInFiledata, err := os.ReadFile(filePath)
		if err != nil {
			return "", errors.New("error: Could not read the file")
		}

		// Set the file data
		builder.FileData = readInFiledata

	} else {
		// Set the file data
		builder.FileData = yamlData
	}

	// Determine the file type
	builder.DetermineFileType()

	// Build the struct from the data
	structOutput := builder.BuildStructFromData()

	// Return the struct output if it was built successfully
	if structOutput != "" {
		return structOutput, nil
	} else {
		return "", errors.New("error: Could not build the struct")
	}
}
