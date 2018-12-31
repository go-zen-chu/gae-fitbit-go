package gcalauth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGcalauth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gcalauth Suite")
}
