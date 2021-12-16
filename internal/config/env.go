package config

import (
	"os"
	"strings"
)

func readEnv(dist map[string]interface{}) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		dist[strings.ToLower(pair[0])] = pair[1]
	}
}
