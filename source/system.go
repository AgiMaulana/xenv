package source

import (
	"os"
	"strings"
)

type SystemSource struct {
}

func (s *SystemSource) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (s *SystemSource) All() map[string]string {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		env[key] = value
	}
	return env
}
