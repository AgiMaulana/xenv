package source

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type FileSource struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewFileSource(filePath string) (*FileSource, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if found {
			data[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &FileSource{
		data: data,
	}, nil
}

// Lookup satisfies the Source interface contract thread-safely.
func (f *FileSource) Lookup(key string) (string, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	val, found := f.data[key]
	return val, found
}

func (f *FileSource) All() map[string]string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	all := make(map[string]string, len(f.data))
	for key, value := range f.data {
		all[key] = value
	}
	return all
}
