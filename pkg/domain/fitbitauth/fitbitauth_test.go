package fitbitauth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fitbitauth", func() {

	var (
		c     *gomock.Controller
		mfbaf *MockFactory
		mfbas *MockStore
		mfhc  *MockFitbitHTTPClient
		ftp   *FitbitTokenParams
		fap   *FitbitAuthParams
		fah   FitbitAuthHandler
		ft    *FitbitTokens
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mfbaf = NewMockFactory(c)
		mfbas = NewMockStore(c)
		mfhc = NewMockFitbitHTTPClient(c)
		fap = &FitbitAuthParams{
			ClientID:     "fb-client-id",
			Scope:        "sleep activity",
			RedirectURI:  "http://127.0.0.1:8080/v1/fitbitstoretoken",
			ResponseType: "code",
			Expires:      "2592000", // 1 week
		}
		ftp = &FitbitTokenParams{
			ClientID:     "fb-client-id",
			ClientSecret: "fb-client-secret",
			GrantType:    "authorization_code",
			RedirectURI:  "http://127.0.0.1:8080/v1/fitbitstoretoken",
		}
		fah = NewFitbitAuthHandler(mfbaf, fap, ftp)
		ft = &FitbitTokens{
			AccessToken:  "access-token",
			ExpiresIn:    28800,
			RefreshToken: "refresh-token",
			Scope:        "activity,sleep",
			TokenType:    "refresh",
			UserID:       "user-id-007",
		}
	})

	Describe("createFitbitAuthURL", func() {
		It("should make fitbit auth url", func() {
			u := createFitbitAuthURL(fap)
			expect := "https://www.fitbit.com/oauth2/authorize?client_id=fb-client-id&expires_in=2592000&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080%2Fv1%2Ffitbitstoretoken&response_type=code&scope=sleep+activity"
			Expect(u.String()).To(Equal(expect))
		})
	})

	Describe("HandleFitbitAuthCode", func() {
		It("should handle fitbit auth code", func() {
			mfbas.EXPECT().WriteFitbitTokens(gomock.Any()).Return(nil)
			mfhc.EXPECT().GetFitbitToken(gomock.Any(), gomock.Any()).Return(ft, nil)

			mfbaf.EXPECT().FileStore().Return(mfbas, nil)
			mfbaf.EXPECT().FitbitHTTPClient().Return(mfhc)
			ts := httptest.NewServer(http.HandlerFunc(fah.HandleFitbitAuthCode))
			res, err := http.Get(ts.URL + "?code=fitbit_auth_code")
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(string(bodyBytes)).To(Equal("OK"))
		})
	})
})
