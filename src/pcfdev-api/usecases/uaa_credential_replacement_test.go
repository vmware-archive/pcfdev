package usecases_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"pcfdev-api/usecases"
)

var _ = Describe("Usecase: UAA Credential Replacement", func() {

	var u *usecases.UaaCredentialReplacement
	goodUaaConfig := `
---
name: uaa
scim:
  some-key: true
  users:
    - admin|admin|other-value
`
	expectedSecuredConfig := `name: uaa
scim:
  some-key: true
  users:
  - admin|some-password|other-value
`

	uaaConfigWithoutScimKey := "name: uaa"
	uaaConfigWithoutUsersKey := "scim: {}"
	uaaConfigWithoutAdminUser := `
---
scim:
  users:
    - other-user|password|other-value
`

	BeforeEach(func() {
		u = &usecases.UaaCredentialReplacement{}
	})

	Context("when the config file is in the expected state", func() {
		It("replaces the admin password", func() {
			securedConfig, err := u.ReplaceUaaConfigAdminCredentials(goodUaaConfig, "some-password")
			Expect(err).NotTo(HaveOccurred())
			Expect(securedConfig).To(Equal(expectedSecuredConfig))
		})
	})

	Context("when the config file is not valid yaml", func() {
		It("returns an error", func() {
			uaaConfig := "some-bad-yaml"
			_, err := u.ReplaceUaaConfigAdminCredentials(uaaConfig, "some-password")
			Expect(err).To(MatchError("failed to parse UAA config file"))
		})
	})

	Context("When the scim key is missing", func() {
		It("returns an error", func() {
			_, err := u.ReplaceUaaConfigAdminCredentials(uaaConfigWithoutScimKey, "some-password")
			Expect(err).To(MatchError("failed to parse UAA config file"))
		})
	})

	Context("When the scim.users key is missing", func() {
		It("returns an error", func() {
			_, err := u.ReplaceUaaConfigAdminCredentials(uaaConfigWithoutUsersKey, "some-password")
			Expect(err).To(MatchError("failed to parse UAA config file"))
		})
	})

	Context("when the admin credentials are not in scim.users", func() {
		It("returns an error", func() {
			_, err := u.ReplaceUaaConfigAdminCredentials(uaaConfigWithoutAdminUser, "some-password")
			Expect(err).To(MatchError("failed to parse UAA config file"))
		})
	})
})
