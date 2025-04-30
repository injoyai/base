package coding

import "encoding/xml"

type Xml struct{}

func (this *Xml) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

func (this *Xml) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

var (
	XmlMarshal   = xml.Marshal
	XmlUnmarshal = xml.Unmarshal
)
