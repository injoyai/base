package coding

import "github.com/pelletier/go-toml/v2"

type Toml struct{}

func (this *Toml) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func (this *Toml) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

var (
	TomlMarshal   = toml.Marshal
	TomlUnmarshal = toml.Unmarshal
)
