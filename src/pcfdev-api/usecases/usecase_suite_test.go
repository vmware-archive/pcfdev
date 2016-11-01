package usecases_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUsecase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usecase Suite")
}

