package runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPCFDevApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PCF Dev Agent Runner Suite")
}
