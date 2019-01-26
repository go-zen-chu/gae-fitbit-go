package gcalauth

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("gcalauth", func() {
	var (
		c     *gomock.Controller
		mf *MockFactory
		ms *MockStore
		moc *MockOAuthClient
		gah   GCalAuthHandler
		config *oauth2.Config
		token    *oauth2.Token
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mf = NewMockFactory(c)
		ms = NewMockStore(c)
		moc = NewMockOAuthClient(c)

		config = &oauth2.Config{
			ClientID: "gcal-client-id",
			ClientSecret: "gcal-client-secret",
			Endpoint: fitbit.Endpoint,
			RedirectURL: "http://127.0.0.1:8080/v1/gcalstoretoken",
			Scopes: []string { calendar.CalendarEventsScope },
		}
		token = &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "refresh",
			Expiry: time.Now().Add(8 * time.Hour),
		}
		mf.EXPECT().FileStore().Return(ms, nil)
		mf.EXPECT().OAuthClient(config).Return(moc)

		gah = NewGCalAuthHandler(mf, config)
	})

	Describe("HandleGCalAuthCode", func() {
		It("should handle gcal auth code", func() {
			moc.EXPECT().Exchange(gomock.Any()).Return(token, nil)
			ms.EXPECT().WriteGCalToken(gomock.Any()).Return(nil)

			ts := httptest.NewServer(http.HandlerFunc(gah.HandleGCalAuthCode))
			res, err := http.Get(ts.URL + "?code=auth_code")
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(string(bodyBytes)).To(Equal("OK"))
		})
	})
})
