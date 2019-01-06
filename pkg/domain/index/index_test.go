package index_test

import (
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Index", func() {
	Describe("HandleIndex", func() {

		var (
			ih index.IndexHandler
			ts *httptest.Server
		)

		BeforeEach(func () {
			ih = index.NewIndexHandler()
			ts = httptest.NewServer(http.HandlerFunc(ih.HandleIndex))
		})

		AfterEach(func(){
			defer ts.Close()
		})

		It("should serve index OK", func() {
			res, err := http.Get(ts.URL)
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(bodyBytes)).To(Equal("OK"))
		})
	})
})


