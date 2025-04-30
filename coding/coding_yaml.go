package coding

import "gopkg.in/yaml.v3"

type Yaml struct{}

func (this *Yaml) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (this *Yaml) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

var (
	YamlMarshal   = yaml.Marshal
	YamlUnmarshal = yaml.Unmarshal
)
