package credentials_generator

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// TRANSPORT_CREDENTIALS defines what security is passed to credentials_generator.Dial and (can be overwritten)
	// you can provide your own, or use credentials_generator.TRANSPORT_CREDENTIALS = credentials_generator.InsecureTransportCredentials
	// if you want no transport credentials (do not use this in production as nothing will get encrypted).
	TRANSPORT_CREDENTIALS TransportCredentials
	// API_ADDR specifies which url we should fetch certificate from if TRANSPORT_CREDENTIALS is not set
	API_ADDR = "api.softcorp.io:443"
)

type TransportCredentials interface {
	GetTransportCredentials() (credentials.TransportCredentials, error)
}

type InsecureTransportCredentials struct{}

func (tc *InsecureTransportCredentials) GetTransportCredentials() (credentials.TransportCredentials, error) {
	return insecure.NewCredentials(), nil
}

type defaultTransportCredentials struct{}

func New() (TransportCredentials, error) {
	if TRANSPORT_CREDENTIALS != nil {
		return TRANSPORT_CREDENTIALS, nil
	}
	return &defaultTransportCredentials{}, nil
}

func (tc *defaultTransportCredentials) GetTransportCredentials() (credentials.TransportCredentials, error) {
	conn, err := tls.Dial("tcp", API_ADDR, &tls.Config{})
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
