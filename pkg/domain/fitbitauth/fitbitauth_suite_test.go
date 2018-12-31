package fitbitauth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFitbitauth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fitbitauth Suite")
}
