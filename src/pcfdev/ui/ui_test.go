package ui_test

import (
	"bytes"
	"errors"
	"pcfdev/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UI", func() {
	Describe("#PrintHelpText", func() {
		var (
			u      *ui.UI
			stdout *bytes.Buffer
		)

		BeforeEach(func() {
			stdout = new(bytes.Buffer)

			u = &ui.UI{
				Stdout: stdout,
			}
		})

		It("should print the help text", func() {
			expectedHelpText := ` _______  _______  _______    ______   _______  __   __
|       ||       ||       |  |      | |       ||  | |  |
|    _  ||       ||    ___|  |  _    ||    ___||  |_|  |
|   |_| ||       ||   |___   | | |   ||   |___ |       |
|    ___||      _||    ___|  | |_|   ||    ___||       |
|   |    |     |_ |   |      |       ||   |___  |     |
|___|    |_______||___|      |______| |_______|  |___|
is now running.
To begin using PCF Dev, please run:
    cf login -a https://api.some-domain --skip-ssl-validation
Apps Manager URL: https://some-domain
Admin user => Email: admin / Password: admin
Regular user => Email: user / Password: pass
`

			Expect(u.PrintHelpText("some-domain")).To(Succeed())
			Expect(stdout.String()).To(Equal(expectedHelpText))
		})

		Context("when there is an error printing the help text", func() {
			It("should return the error", func() {
				u.Stdout = &brokenWriter{}

				Expect(u.PrintHelpText("some-domain")).To(MatchError("some-error"))
			})
		})
	})
})

type brokenWriter struct{}

func (*brokenWriter) Write(b []byte) (int, error) {
	return 0, errors.New("some-error")
}
