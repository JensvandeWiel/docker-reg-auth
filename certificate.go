package registry

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/docker/libtrust"
)

// LoadCertificateAndKey loads a certificate and key from the given paths and returns the public and private keys. This expects x509 certificates.
func LoadCertificateAndKey(crtPath, keyPath string) (libtrust.PublicKey, libtrust.PrivateKey, error) {
	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, nil, err
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := libtrust.FromCryptoPublicKey(x509Cert.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := libtrust.FromCryptoPrivateKey(cert.PrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return pubKey, privKey, nil
}
