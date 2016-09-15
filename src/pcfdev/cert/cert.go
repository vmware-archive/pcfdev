package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type Cert struct{}

func (c *Cert) GenerateCert(domain string) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now().Add(-24 * time.Hour)
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	template := x509.Certificate{
		DNSNames: []string{
			domain,
			"*." + domain,
			"*.login." + domain,
			"*.uaa." + domain,
		},
		EmailAddresses: []string{"pcfdev-eng@pivotal.io"},
		NotBefore:      notBefore,
		NotAfter:       notAfter,
		SerialNumber:   big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         domain,
			Country:            []string{"US"},
			Locality:           []string{"New York"},
			Organization:       []string{"Cloud Foundry"},
			OrganizationalUnit: []string{"PCF Dev"},
			Province:           []string{"New York"},
		},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	encodedCertificate := new(bytes.Buffer)
	if err := pem.Encode(encodedCertificate, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
		return nil, nil, err
	}

	encodedPrivateKey := new(bytes.Buffer)
	if err := pem.Encode(encodedPrivateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return nil, nil, err
	}

	return encodedCertificate.Bytes(), encodedPrivateKey.Bytes(), nil
}
