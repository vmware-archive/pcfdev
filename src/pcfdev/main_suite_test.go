package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPCFDev(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PCF Dev Main Suite")
}
