package command_test

import (
	"github.com/go-zen-chu/gae-fitbit-go/pkg/application/command"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"


)

var _ = Describe("command", func() {
	var (
		c *gomock.Controller
		mih *index.MockIndexHandler
		mfaf *dfba.MockFactory
		mfah *dfba.MockFitbitAuthHandler
		mgaf *dga.MockFactory
		mgah *dga.MockGCalAuthHandler
		mf2gf *df2g.MockFactory
		mf2gs *df2g.MockService
		mhs *command.MockHttpServer
	)

	BeforeEach(func(){
		c = gomock.NewController(GinkgoT())
		mih = index.NewMockIndexHandler(c)
		mfaf = dfba.NewMockFactory(c)
		mfah = dfba.NewMockFitbitAuthHandler(c)
		mgaf = dga.NewMockFactory(c)
		mgah = dga.NewMockGCalAuthHandler(c)
		mf2gf = df2g.NewMockFactory(c)
		mf2gs = df2g.NewMockService(c)
		mhs = command.NewMockHttpServer(c)
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("NewCommand", func(){
		It("should make new command", func() {
			cmd := command.NewCommand(mih, mfaf, mgaf, mf2gf, mhs)
			Expect(cmd).ShouldNot(BeNil())
		})
	})

	Describe("Run", func() {
		var cmd command.Command

		BeforeEach(func() {
			mfaf.EXPECT().FitbitAuthHandler(gomock.Any(), gomock.Any()).Return(mfah)
			mgaf.EXPECT().GCalAuthHandler(gomock.Any()).Return(mgah)
			mf2gf.EXPECT().Service(gomock.Any()).Return(mf2gs)

			mhs.EXPECT().Run("9091").Return(nil)

			cmd = command.NewCommand(mih, mfaf, mgaf, mf2gf, mhs)
		})

		It("should run", func() {
			args := os.Args
			os.Args = []string{"gae-fitbit-go",
				"--port", "9091",
				"--fb-client-id", "fb-client-id",
				"--fb-client-secret", "fb-client-secret",
				"--gcal-sleep-cal-id", "gcal-sleep-cal-id",
				"--gcal-activity-cal-id", "gcal-activity-cal-id",
				"--gcal-client-id", "gcal-client-id",
				"--gcal-client-secret", "gcal-client-secret",
			}
			err := cmd.Run()
			Expect(err).To(BeNil())
			os.Args = args
		})
	})

})
