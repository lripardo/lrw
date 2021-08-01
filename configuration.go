package lrw

type Configuration interface {
	String(ConfigurationKey) string
	Uint(ConfigurationKey) uint
	Bool(ConfigurationKey) bool
	Int(ConfigurationKey) int
	Int64(ConfigurationKey) int64
	Strings(ConfigurationKey) []string
}
