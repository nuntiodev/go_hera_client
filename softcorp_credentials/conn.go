package softcorp_credentials

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type TransportCredentials interface {
	GetTransportCredentials(api string) (credentials.TransportCredentials, error)
}
type InsecureTransportCredentials struct{}

func (tc *InsecureTransportCredentials) GetTransportCredentials() (credentials.TransportCredentials, error) {
	return insecure.NewCredentials(), nil
}

type defaultTransportCredentials struct{}

func New(transportCredentials TransportCredentials) (TransportCredentials, error) {
	if transportCredentials != nil {
		return transportCredentials, nil
	}
	return &defaultTransportCredentials{}, nil
}

func (tc *defaultTransportCredentials) GetTransportCredentials(api string) (credentials.TransportCredentials, error) {
	conn, err := tls.Dial("tcp", api, &tls.Config{})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var b bytes.Buffer
	for _, cert := range conn.ConnectionState().PeerCertificates {
		err := pem.Encode(&b, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return nil, err
		}
	}
	b.String()
	clientTLSCertPool := x509.NewCertPool()
	clientTLSCertPool.AppendCertsFromPEM([]byte(b.String()))
	return credentials.NewClientTLSFromCert(clientTLSCertPool, ""), nil
}
