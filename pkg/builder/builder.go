package builder

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed version
var VersionEmbed embed.FS

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
	StructName string
	// The path to the file
	FilePath string
	// fileData is the data from the file
	FileData []byte
	// The file type of the provided file (JSON or YAML)
	FileType string
	// The data from the file
	mappedFileData map[interface{}]interface{}
	// structOutput is the output of the struct
	StructOutput string
)

// DetermineFileType determines the file type of the provided file
func DetermineFileType() {

	// If the file data is provided, don't read the file
	if FilePath != "" {
		// Read the file
		readInFiledata, err := os.ReadFile(FilePath)
		if err != nil {
			fmt.Println("Error: Could not read the file.")
			os.Exit(1)
		}
		FileData = readInFiledata
	}

	// Try to unmarshal the data into a map[interface{}]interface{}
	var data map[interface{}]interface{}
	jsonErr := json.Unmarshal(FileData, &data)
	yamlErr := yaml.Unmarshal(FileData, &data)

	// If JSON unmarshalling was successful, set the file type to JSON
	if jsonErr == nil {
		FileType = ftJSON
		mappedFileData = data
		return
	}

	// If YAML unmarshalling was successful, set the file type to YAML
	if yamlErr == nil {
		FileType = ftYAML
		mappedFileData = data
		return
	}

	// Otherwise, the file type is unsupported
	fmt.Println("Unsupported file type.")
	os.Exit(1)
}

// buildStructFromData builds a struct from the YAML data
func BuildStructFromData() string {
	// Create a struct from the YAML data
	StructOutput = "// Generated using github.com/mitchs-dev/build-struct@" + GetVersion() + "\n"
	StructOutput += "type " + StructName + " struct {\n"
	StructOutput += structBuilder(mappedFileData, "\t")
	StructOutput += "}"

	return StructOutput
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
		case map[interface{}]interface{}:
			fieldType = "struct {\n" + structBuilder(v, prefix+"\t") + prefix + "}"
		case []interface{}:
			if len(v) > 0 {
				// Check if the first element is a map
				if _, ok := v[0].(map[interface{}]interface{}); ok {
					// If it's a map, generate a slice of structs
					fieldType = "[]struct {\n" + structBuilder(v[0].(map[interface{}]interface{}), prefix+"\t") + prefix + "}"
				} else {
					// If it's not a map, use the existing logic
					types := make(map[string]bool)
					for _, elem := range v {
						types[getType(elem)] = true
					}
					if len(types) == 1 {
						for t := range types {
							fieldType = "[]" + t
						}
					} else {
						fieldType = "[]interface{}"
					}
				}
			} else {
				fieldType = "[]interface{}"
			}
		default:
			fieldType = getType(v)
		}
		keyStr := fmt.Sprintf("%v", key)
		structFields += prefix + fmt.Sprintf("%s %s `"+FileType+":\"%s\"`\n", cases.Title(language.English).String(keyStr), fieldType, key)
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

// GetVersion returns the version of the application
func GetVersion() string {
	versionData, err := VersionEmbed.ReadFile("version")
	if err != nil {
		return "Unknown"
	}
	return string(versionData)
}
