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

func (c *Cert) GenerateCerts(domain string) ([]byte, []byte, []byte, []byte, error) {
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	encodedCAPrivateKey := new(bytes.Buffer)
	if err := pem.Encode(encodedCAPrivateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey)}); err != nil {
		return nil, nil, nil, nil, err
	}

	caTemplate := c.generateTemplate(domain, true)
	encodedCACertificate, err := c.generateCert(caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	encodedPrivateKey := new(bytes.Buffer)
	if err := pem.Encode(encodedPrivateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return nil, nil, nil, nil, err
	}

	template := c.generateTemplate(domain, false)
	encodedCertificate, err := c.generateCert(template, caTemplate, &privateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return encodedCertificate, encodedPrivateKey.Bytes(), encodedCACertificate, encodedCAPrivateKey.Bytes(), nil
}

func (c *Cert) generateCert(template, parentTemplate *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) ([]byte, error) {
	certificate, err := x509.CreateCertificate(rand.Reader, template, parentTemplate, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	encodedCertificate := new(bytes.Buffer)
	if err := pem.Encode(encodedCertificate, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
		return nil, err
	}

	return encodedCertificate.Bytes(), nil
}

func (c *Cert) generateTemplate(domain string, isCA bool) *x509.Certificate {
	notBefore := time.Now().Add(-24 * time.Hour)
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	template := &x509.Certificate{
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
			CommonName:         domain + " " + time.Now().Format(time.RFC3339),
			Country:            []string{"US"},
			Locality:           []string{"New York"},
			Organization:       []string{"Cloud Foundry"},
			OrganizationalUnit: []string{"PCF Dev"},
			Province:           []string{"New York"},
		},
	}

	if isCA {
		template.Subject.Organization = []string{"Cloud Foundry CA"}
		template.BasicConstraintsValid = true
		template.IsCA = true
	}

	return template
}
