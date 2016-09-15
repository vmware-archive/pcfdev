package cert_test

import (
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"pcfdev/cert"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cert", func() {
	Describe("#GenerateCert", func() {
		var c *cert.Cert

		BeforeEach(func() {
			c = &cert.Cert{}
		})

		It("should generate a self-signed certificate and private key", func() {
			certificateBytes, privateKeyBytes, err := c.GenerateCert("some-domain")
			Expect(err).NotTo(HaveOccurred())

			yesterday := time.Now().Add(-24 * time.Hour)
			tenYearsFromYesterday := yesterday.Add(10 * 365 * 24 * time.Hour)

			certificate := parseCertificate(certificateBytes, privateKeyBytes)

			Expect(certificate.DNSNames).To(Equal([]string{
				"some-domain",
				"*.some-domain",
				"*.login.some-domain",
				"*.uaa.some-domain",
			}))
			Expect(certificate.EmailAddresses).To(Equal([]string{"pcfdev-eng@pivotal.io"}))
			Expect(certificate.Issuer).To(Equal(certificate.Subject))
			Expect(certificate.NotBefore).To(BeTemporally("~", yesterday, time.Minute))
			Expect(certificate.NotAfter).To(BeTemporally("~", tenYearsFromYesterday, time.Minute))
			Expect(certificate.SerialNumber).To(Equal(big.NewInt(1)))
			Expect(certificate.Subject.CommonName).To(Equal("some-domain"))
			Expect(certificate.Subject.Country).To(Equal([]string{"US"}))
			Expect(certificate.Subject.Locality).To(Equal([]string{"New York"}))
			Expect(certificate.Subject.Organization).To(Equal([]string{"Cloud Foundry"}))
			Expect(certificate.Subject.OrganizationalUnit).To(Equal([]string{"PCF Dev"}))
			Expect(certificate.Subject.Province).To(Equal([]string{"New York"}))
			Expect(certificate.IsCA).To(BeTrue())
		})
	})
})

func parseCertificate(certificateBytes []byte, privateKeyBytes []byte) *x509.Certificate {
	certificate, err := tls.X509KeyPair(certificateBytes, privateKeyBytes)
	Expect(err).NotTo(HaveOccurred())
	certificates, err := x509.ParseCertificates(certificate.Certificate[0])
	Expect(err).NotTo(HaveOccurred())
	return certificates[0]
}
