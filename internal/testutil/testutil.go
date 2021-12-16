package testutil

import (
	"os"
	"regexp"
)

// SetRootPath sets root folder for tests
func SetRootPath() {
	currentWorkDirectory, _ := os.Getwd()

	rootPath := regexp.
		MustCompile(`(.+)/internal/.+`).
		ReplaceAllString(currentWorkDirectory, `$1`)

	err := os.Chdir(rootPath)

	if err != nil {
		panic(err)
	}
}
