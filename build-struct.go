/*
This is a tool designed to build a struct from a given JSON or YAML configuration file.

Usage:

	build-struct <struct-name> <file>
*/
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Static variables
var (
	// Variable for the file type - JSON
	ftJSON = "json"
	// Variable for the file type - YAML
	ftYAML = "yaml"
)

// Dynamic variables
var (
	// structName is the name of the struct
	structName string
	// The path to the file
	filePath string
	// The file type of the provided file (JSON or YAML)
	fileType string
	// The data from the file
	mappedFileData map[interface{}]interface{}
	// structOutput is the output of the struct
	structOutput string
)

func init() {

	// Make sure that a file is provided as an argument
	if len(os.Args) < 3 || len(os.Args) > 3 {
		fmt.Println("Usage: build-struct <struct-name> <file>")
		os.Exit(1)
	}

	// Set the struct name
	structName = os.Args[1]

	// Set the file path
	filePath = os.Args[2]

	// Make sure that the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Error: The provided file does not exist.")
		os.Exit(1)
	}

	// Make sure that the file is not empty
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Error: The file is empty.")
		os.Exit(1)
	}

	// Make sure that it's a file and not a directory
	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo.IsDir() {
		fmt.Println("Error: The provided path is a directory.")
		os.Exit(1)
	}

}

func main() {

	determineFileType()
	structOutput := buildStructFromData()

	// Return the struct output if it was built successfully
	if structOutput != "" {
		fmt.Println(structOutput)
		os.Exit(0)
	} else {
		fmt.Println("Error: Could not build the struct.")
		os.Exit(1)
	}

}

// determineFileType determines the file type of the provided file
func determineFileType() {
	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error: Could not read the file.")
		os.Exit(1)
	}

	// Try to unmarshal the data into a map[interface{}]interface{}
	var data map[interface{}]interface{}
	jsonErr := json.Unmarshal(fileData, &data)
	yamlErr := yaml.Unmarshal(fileData, &data)

	// If JSON unmarshalling was successful, set the file type to JSON
	if jsonErr == nil {
		fileType = ftJSON
		mappedFileData = data
		return
	}

	// If YAML unmarshalling was successful, set the file type to YAML
	if yamlErr == nil {
		fileType = ftYAML
		mappedFileData = data
		return
	}

	// Otherwise, the file type is unsupported
	fmt.Println("Unsupported file type.")
	os.Exit(1)
}

// buildStructFromData builds a struct from the YAML data
func buildStructFromData() string {
	// Create a struct from the YAML data
	structOutput = "type " + structName + " struct {\n"
	structOutput += structBuilder(mappedFileData, "\t")
	structOutput += "}"

	return structOutput
}

// structBuilder builds a struct from the provided data
func structBuilder(data map[interface{}]interface{}, prefix string) string {

	// Initialize the struct fields
	var structFields string

	// Iterate over the data
	for key, value := range data {

		// Get the type of the value
		var fieldType string

		// Check the type of the value
		switch v := value.(type) {
		// If the value is a map, recursively build a struct for it
		case map[interface{}]interface{}:
			// If the value is a map, recursively build a struct for it
			fieldType = "struct {\n" + structBuilder(v, prefix+"\t") + prefix + "}"
		// If the value is a slice, check each element's type
		case []interface{}:
			// If the value is a slice, check each element's type
			for _, elem := range v {
				// Check the type of the element
				switch elem := elem.(type) {
				// If the element is a map, build a struct for it
				case map[interface{}]interface{}:
					// If the element is a map, build a struct for it
					fieldType = "[]struct {\n" + structBuilder(elem, prefix+"\t") + prefix + "}"
				// If the element is not a map, get its type
				default:
					// If the element is not a map, get its type
					fieldType = "[]" + getType(elem)
				}
			}
			if fieldType == "" {
				// If the slice is empty, default to []interface{}
				fieldType = "[]interface{}"
			}
		// If the value is not a map or a slice, get its type
		default:
			fieldType = getType(v)
		}
		keyStr := fmt.Sprintf("%v", key)
		structFields += prefix + fmt.Sprintf("%s %s `"+fileType+":\"%s\"`\n", strings.Title(keyStr), fieldType, key)
	}

	return structFields
}

// getType returns the type of the provided value
func getType(v interface{}) string {
	switch v := v.(type) {
	case string:
		return "string"
	case float64:
		return "float64"
	case int:
		if math.Abs(float64(v)) > math.MaxInt32 {
			return "int64"
		} else {
			return "int"
		}
	case bool:
		return "bool"
	case []interface{}:
		return "[]interface{}"
	case map[interface{}]interface{}:
		return "struct {\n" + structBuilder(v, "\t") + "}\n"

	default:
		return "interface{}"
	}
}
