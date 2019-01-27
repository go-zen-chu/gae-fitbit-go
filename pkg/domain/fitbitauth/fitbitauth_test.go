package fitbitauth

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("fitbitauth", func() {
	var (
		c      *gomock.Controller
		mf     *MockFactory
		ms     *MockStore
		moc    *MockOAuthClient
		fah    FitbitAuthHandler
		config *oauth2.Config
		token  *oauth2.Token
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mf = NewMockFactory(c)
		ms = NewMockStore(c)
		moc = NewMockOAuthClient(c)

		config = &oauth2.Config{
			ClientID:     "fb-client-id",
			ClientSecret: "fb-client-secret",
			Endpoint:     fitbit.Endpoint,
			RedirectURL:  "http://127.0.0.1:8080/v1/fitbitstoretoken",
			Scopes:       []string{"sleep", "activity"},
		}
		token = &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "refresh",
			Expiry:       time.Now().Add(8 * time.Hour),
		}
		mf.EXPECT().FileStore().Return(ms, nil)
		mf.EXPECT().OAuthClient(config).Return(moc)

		fah = NewFitbitAuthHandler(mf, config)
	})

	Describe("HandleFitbitAuthCode", func() {
		It("should handle fitbit auth code", func() {
			moc.EXPECT().Exchange(gomock.Any()).Return(token, nil)
			ms.EXPECT().WriteFitbitToken(gomock.Any()).Return(nil)

			ts := httptest.NewServer(http.HandlerFunc(fah.HandleFitbitAuthCode))
			res, err := http.Get(ts.URL + "?code=auth_code")
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(string(bodyBytes)).To(Equal("OK"))
		})
	})
})
