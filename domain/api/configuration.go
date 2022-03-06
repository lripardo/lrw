package api

type Configuration interface {
	String(key Key) string
	Uint(key Key) uint
	Bool(key Key) bool
	Strings(key Key) []string
	Int64(key Key) int64
	Int(key Key) int
}

type Key struct {
	name       string
	validation string
	value      string
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Validation() string {
	return k.validation
}

func (k *Key) Value() string {
	return k.value
}

func NewKey(name, validation, value string) Key {
	return Key{
		name:       name,
		validation: validation,
		value:      value,
	}
}
