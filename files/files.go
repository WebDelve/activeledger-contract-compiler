package files

import (
	"fmt"
	"os"

	"github.com/WebDelve/activeledger-contract-compiler/helper"
)

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		helper.Error(err, fmt.Sprintf("Error reading file with path %s", path))
	}

	return data
}

func WriteFile(path string, bData []byte) {
	if err := os.WriteFile(path, bData, 0644); err != nil {
		helper.Error(err, fmt.Sprintf("Error writing data to file \"%s\"", path))
	}
}
