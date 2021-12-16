package config

import (
	"testing"

	testutil2 "github.com/ArtemGretsov/golang-rabbitmq-template/internal/testutil"
)

func init() {
	testutil2.SetRootPath()
}

func Test_ReadJSONConfig(t *testing.T) {
	t.Parallel()

	dist := make(map[string]interface{})
	err := readJSONConfig(DefaultPath, dist, true)

	if err != nil {
		t.Error(err)
	}
}
