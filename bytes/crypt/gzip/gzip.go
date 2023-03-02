package gzip

import (
	"bytes"
	"compress/gzip"
)

// EncodeGzip 压缩字节
func EncodeGzip(input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(buf)
	_, err := gzipWriter.Write(input)
	gzipWriter.Close()
	return buf.Bytes(), err
}

// DecodeGzip 解压字节
func DecodeGzip(input []byte) ([]byte, error) {
	reader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(gzipReader)
	return buf.Bytes(), err
}
