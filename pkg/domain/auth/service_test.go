package auth_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	. "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/auth"
)

var _ = Describe("Service", func() {
	var (
		c    *gomock.Controller
		ms   *MockStore
		mfoc *MockOAuthClient
		mgoc *MockOAuthClient
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		ms = NewMockStore(c)
		mfoc = NewMockOAuthClient(c)
		mgoc = NewMockOAuthClient(c)
	})

	Describe("HandleFitbitAuthCode", func() {
		var (
			svc   Service
			code  = "some nice auth code here"
			token *oauth2.Token
		)

		BeforeEach(func() {
			svc = NewService(ms, mfoc, mgoc)
			token = &oauth2.Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}
			mfoc.EXPECT().Exchange(code).Return(token, nil)
			ms.EXPECT().WriteFitbitToken(gomock.Any()).Return(nil)
		})

		It("should handle fitbit auth code", func() {
			err := svc.HandleFitbitAuthCode(code)
			Expect(err).To(BeNil())
		})
	})

	Describe("HandleGCalAuthCode", func() {
		var (
			svc   Service
			code  = "some nice auth code here"
			token *oauth2.Token
		)

		BeforeEach(func() {
			svc = NewService(ms, mfoc, mgoc)
			token = &oauth2.Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}
			mgoc.EXPECT().Exchange(code).Return(token, nil)
			ms.EXPECT().WriteGCalToken(gomock.Any()).Return(nil)
		})

		It("should handle gcal auth code", func() {
			err := svc.HandleGCalAuthCode(code)
			Expect(err).To(BeNil())
		})
	})
})
