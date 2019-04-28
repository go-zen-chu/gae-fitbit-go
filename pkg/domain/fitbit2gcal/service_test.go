package fitbit2gcal_test

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
)

var _ = Describe("Service", func() {
	var (
		c   *gomock.Controller
		mfc *MockFitbitClient
		mgc *MockGCalClient
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mfc = NewMockFitbitClient(c)
		mgc = NewMockGCalClient(c)
	})

	Describe("Fitbit2GCal", func() {
		var (
			svc      Service
			fromDate time.Time
			toDate   time.Time
		)

		BeforeEach(func() {
			svc = NewService(mfc, mgc)
		})

		Context("if fromDate > toDate", func() {
			It("should be error", func() {
				fromDate = time.Now()
				toDate = fromDate.Add(time.Hour * -48)
				err := svc.Fitbit2GCal(fromDate, toDate)
				Expect(err).NotTo(BeNil())
			})
		})

		//TODO: TBD
		// Context("if proper params are given", func() {
		// 	It("should run", func(){
		// 		mfc.EXPECT().Exchange(code).Return(token, nil)
		// 		ms.EXPECT().WriteFitbitToken(gomock.Any()).Return(nil)
		//
		// 	})
		// })
	})
})
