package fitbit2gcal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFitbit2gcal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fitbit2gcal Suite")
}
