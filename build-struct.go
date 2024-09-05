/*
This is a tool designed to build a struct from a given JSON or YAML configuration file.

Usage:

	build-struct <struct-name> <file>
*/
package main

import (
	"fmt"
	"os"

	"github.com/mitchs-dev/build-struct/pkg/builder"
)

func init() {

	// Make sure that a file is provided as an argument
	if len(os.Args) < 3 || len(os.Args) > 3 {
		fmt.Println("Usage: build-struct <struct-name> <file>")
		os.Exit(1)
	}

	// Set the struct name
	builder.StructName = os.Args[1]

	// Set the file path
	builder.FilePath = os.Args[2]

	// Make sure that the file exists
	if _, err := os.Stat(builder.FilePath); os.IsNotExist(err) {
		fmt.Println("Error: The provided file does not exist.")
		os.Exit(1)
	}

	// Make sure that the file is not empty
	if _, err := os.Stat(builder.FilePath); os.IsNotExist(err) {
		fmt.Println("Error: The file is empty.")
		os.Exit(1)
	}

	// Make sure that it's a file and not a directory
	if fileInfo, err := os.Stat(builder.FilePath); err == nil && fileInfo.IsDir() {
		fmt.Println("Error: The provided path is a directory.")
		os.Exit(1)
	}

}

func main() {

	builder.DetermineFileType()
	structOutput := builder.BuildStructFromData()

	// Return the struct output if it was built successfully
	if structOutput != "" {
		fmt.Println(structOutput)
		os.Exit(0)
	} else {
		fmt.Println("Error: Could not build the struct.")
		os.Exit(1)
	}

}
