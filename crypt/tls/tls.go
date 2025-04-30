package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

type Config struct {
	CAFile string //server
	CCFile string // client-crt
	CKFile string // client-key
}

func (this Config) New() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(this.CAFile)
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca)
	clientKeyPair, err := tls.LoadX509KeyPair(this.CCFile, this.CKFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            certPool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientKeyPair},
	}, nil
}
