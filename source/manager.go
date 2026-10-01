package source

type Manager struct {
	sources []Source
}

func NewManager(sources ...Source) *Manager {
	return &Manager{sources: sources}
}

// Lookup iterates through all sources and returns the first match found.
func (m *Manager) Lookup(key string) (string, bool) {
	for _, src := range m.sources {
		if val, found := src.Lookup(key); found {
			return val, true
		}
	}
	return "", false
}

// All merges every source, keeping the value from the first source that
// defines a given key.
func (m *Manager) All() map[string]string {
	all := make(map[string]string)
	for _, src := range m.sources {
		for key, value := range src.All() {
			if _, exists := all[key]; !exists {
				all[key] = value
			}
		}
	}
	return all
}

func (m *Manager) GetString(key string, defaultVal string) string {
	if val, found := m.Lookup(key); found {
		return val
	}
	return defaultVal
}
