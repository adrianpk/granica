package config

// Config manager.
type Config struct {
	namespace string // i.e.: MWZ, MYCVS, APP, etc...
	values    map[string]string
}
