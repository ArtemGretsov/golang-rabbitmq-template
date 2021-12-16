package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

func readJSONConfig(path string, dist map[string]interface{}, isRequiredExist bool) error {
	file, err := os.Open(path)

	if err != nil {
		if isRequiredExist {
			return errors.WithStack(err)
		}

		return nil
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = json.Unmarshal(content, &dist); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
