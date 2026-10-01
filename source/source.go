package source

type Source interface {
	Lookup(key string) (string, bool)
	All() map[string]string
}
