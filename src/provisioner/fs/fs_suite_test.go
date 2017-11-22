package fs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FS Suite")
}
