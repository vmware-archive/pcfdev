package provisioner_test

import (
	"bytes"
	"errors"
	"pcfdev/provisioner"
	"pcfdev/provisioner/mocks"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioner", func() {
	Describe("#Provision", func() {
		var (
			p             *provisioner.Provisioner
			mockCtrl      *gomock.Controller
			mockCert      *mocks.MockCert
			mockCmdRunner *mocks.MockCmdRunner
			mockFS        *mocks.MockFS
			mockUI        *mocks.MockUI
			firstCommand  *mocks.MockCommand
			secondCommand *mocks.MockCommand
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockCert = mocks.NewMockCert(mockCtrl)
			mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
			mockFS = mocks.NewMockFS(mockCtrl)
			mockUI = mocks.NewMockUI(mockCtrl)
			firstCommand = mocks.NewMockCommand(mockCtrl)
			secondCommand = mocks.NewMockCommand(mockCtrl)

			p = &provisioner.Provisioner{
				Cert:      mockCert,
				CmdRunner: mockCmdRunner,
				FS:        mockFS,
				UI:        mockUI,
				Commands: []provisioner.Command{
					firstCommand,
					secondCommand,
				},

				Distro: provisioner.DistributionPCF,
			}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("should provision a VM", func() {
			gomock.InOrder(
				mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
				mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
				mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
				mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
				mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
				mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))),
				firstCommand.EXPECT().Run(),
				secondCommand.EXPECT().Run(),
				mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
				mockFS.EXPECT().Write("/run/pcfdev-healthcheck", bytes.NewReader([]byte(""))),
			)

			Expect(p.Provision("some-provision-script-path", "some-domain")).To(Succeed())
		})

		Context("when the distribution is oss", func() {
			It("should not run pcf 'Commands'", func() {
				p.Distro = provisioner.DistributionOSS

				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
					mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))),
					firstCommand.EXPECT().Distro().Return(provisioner.DistributionPCF),
					secondCommand.EXPECT().Distro().Return(provisioner.DistributionOSS),
					secondCommand.EXPECT().Run(),
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
					mockFS.EXPECT().Write("/run/pcfdev-healthcheck", bytes.NewReader([]byte(""))),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(Succeed())
			})
		})

		Context("when there is an error generating certificate", func() {
			It("should return the error", func() {
				mockCert.EXPECT().GenerateCerts("some-domain").Return(nil, nil, nil, nil, errors.New("some-error"))

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error creating the gorouter config directory", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the certificate", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the private key", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error creating the openssl directory", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the CA certificate", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
					mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when a command fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
					mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))),
					firstCommand.EXPECT().Run().Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error running the provision script", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
					mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))),
					firstCommand.EXPECT().Run(),
					secondCommand.EXPECT().Run(),
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the healthcheck file", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCerts("some-domain").Return([]byte("some-cert"), []byte("some-key"), []byte("some-ca-cert"), []byte("some-ca-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockFS.EXPECT().Mkdir("/var/pcfdev/openssl"),
					mockFS.EXPECT().Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader([]byte("some-ca-cert"))),
					firstCommand.EXPECT().Run(),
					secondCommand.EXPECT().Run(),
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
					mockFS.EXPECT().Write("/run/pcfdev-healthcheck", bytes.NewReader([]byte(""))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})
	})
})
