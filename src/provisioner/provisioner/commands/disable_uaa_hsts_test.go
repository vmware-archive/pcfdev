package commands_test

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"provisioner/provisioner/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html/charset"
	"provisioner/provisioner"
)

var (
	tempDir    string
	webXMLPath string
)

var _ = Describe("DisableUAAHSTS", func() {

	var cmd *commands.DisableUAAHSTS

	Describe("#Run", func() {
		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "pcfdev-commands")
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		Context("when HSTS is not disabled in tomcat's web.xml", func() {
			BeforeEach(func() {
				templateXMLpath := filepath.Join("..", "..", "assets", "tomcat-web.xml")
				writeTomcatXML(templateXMLpath)

				cmd = &commands.DisableUAAHSTS{
					WebXMLPath: webXMLPath,
				}
			})

			It("should edit tomcat's web.xml to disable HSTS", func() {
				webXMLData := decodeWebXML(webXMLPath)
				Expect(len(webXMLData.Filters)).To(Equal(0))

				Expect(cmd.Run()).To(Succeed())

				webXMLData = decodeWebXML(webXMLPath)
				Expect(len(webXMLData.Filters)).To(Equal(1))
				Expect(webXMLData.Filters[0].FilterName).To(Equal("httpHeaderSecurity"))
				Expect(webXMLData.Filters[0].FilterClass).To(Equal("org.apache.catalina.filters.HttpHeaderSecurityFilter"))
				Expect(webXMLData.Filters[0].InitParam.ParamName).To(Equal("hstsEnabled"))
				Expect(webXMLData.Filters[0].InitParam.ParamValue).To(Equal("false"))
				Expect(webXMLData.Filters[0].AsyncSupported).To(BeTrue())
			})
		})

		Context("when HSTS is already disabled in tomcat's web.xml", func() {
			BeforeEach(func() {
				templateXMLpath := filepath.Join("..", "..", "assets", "tomcat-web-hsts-disabled.xml")
				writeTomcatXML(templateXMLpath)

				cmd = &commands.DisableUAAHSTS{
					WebXMLPath: webXMLPath,
				}
			})

			It("should not add another filter and keep the other filters", func() {
				webXMLData := decodeWebXML(webXMLPath)
				Expect(len(webXMLData.Filters)).To(Equal(2))
				Expect(webXMLData.Filters[0].FilterName).To(Equal("httpHeaderSecurity"))
				Expect(webXMLData.Filters[0].FilterClass).To(Equal("org.apache.catalina.filters.HttpHeaderSecurityFilter"))
				Expect(webXMLData.Filters[0].InitParam.ParamName).To(Equal("hstsEnabled"))
				Expect(webXMLData.Filters[0].InitParam.ParamValue).To(Equal("false"))
				Expect(webXMLData.Filters[0].AsyncSupported).To(BeTrue())

				Expect(webXMLData.Filters[1].FilterName).To(Equal("some-other-filter"))
				Expect(webXMLData.Filters[1].FilterClass).To(Equal("some-other-company"))
				Expect(webXMLData.Filters[1].InitParam.ParamName).To(Equal("some-param"))
				Expect(webXMLData.Filters[1].InitParam.ParamValue).To(Equal("some-value"))

				Expect(cmd.Run()).To(Succeed())

				webXMLData = decodeWebXML(webXMLPath)
				Expect(len(webXMLData.Filters)).To(Equal(2))
				Expect(webXMLData.Filters[0].FilterName).To(Equal("httpHeaderSecurity"))
				Expect(webXMLData.Filters[0].FilterClass).To(Equal("org.apache.catalina.filters.HttpHeaderSecurityFilter"))
				Expect(webXMLData.Filters[0].InitParam.ParamName).To(Equal("hstsEnabled"))
				Expect(webXMLData.Filters[0].InitParam.ParamValue).To(Equal("false"))
				Expect(webXMLData.Filters[0].AsyncSupported).To(BeTrue())
				Expect(webXMLData.Filters[1].FilterName).To(Equal("some-other-filter"))
				Expect(webXMLData.Filters[1].FilterClass).To(Equal("some-other-company"))
				Expect(webXMLData.Filters[1].InitParam.ParamName).To(Equal("some-param"))
				Expect(webXMLData.Filters[1].InitParam.ParamValue).To(Equal("some-value"))
			})
		})

		Context("when the path to the web.xml does not exist", func() {
			BeforeEach(func() {
				cmd = &commands.DisableUAAHSTS{
					WebXMLPath: "/some/bad/path",
				}
			})

			It("should return an error", func() {
				Expect(cmd.Run()).To(MatchError(ContainSubstring("no such file or directory")))
			})
		})

		Context("when the XML being parsed is invalid", func() {
			BeforeEach(func() {
				templateXMLpath := filepath.Join("..", "..", "assets", "tomcat-web-invalid.xml")
				writeTomcatXML(templateXMLpath)

				cmd = &commands.DisableUAAHSTS{
					WebXMLPath: webXMLPath,
				}
			})

			It("should return an error", func() {
				Expect(cmd.Run()).To(MatchError(ContainSubstring("EOF")))
			})
		})
	})

	Describe("#Distro", func() {
		It("should return 'pcf'", func() {
			Expect(cmd.Distro()).To(Equal(provisioner.DistributionPCF))
		})
	})
})

func writeTomcatXML(path string) {
	webXMLContents, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	webXMLPath = filepath.Join(tempDir, "web.xml")
	Expect(ioutil.WriteFile(webXMLPath, webXMLContents, 0644)).To(Succeed())

}

func decodeWebXML(webXMLPath string) *commands.WebApp {
	var webXMLData commands.WebApp
	webXMLReader, err := os.Open(webXMLPath)
	defer webXMLReader.Close()
	Expect(err).NotTo(HaveOccurred())
	decoder := xml.NewDecoder(webXMLReader)
	decoder.CharsetReader = charset.NewReaderLabel
	Expect(decoder.Decode(&webXMLData)).To(Succeed())
	return &webXMLData
}
