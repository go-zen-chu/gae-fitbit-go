package command_test

import (
	"os"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/application/command"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("command", func() {
	var (
		c   *gomock.Controller
		mhs *command.MockHttpServer
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mhs = command.NewMockHttpServer(c)
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("NewCommand", func() {
		It("should make new command", func() {
			cmd := command.NewCommand(mhs)
			Expect(cmd).ShouldNot(BeNil())
		})
	})

	Describe("Run", func() {
		var cmd command.Command
		BeforeEach(func() {
			mhs.EXPECT().HandleFunc(gomock.Any(), gomock.Any()).AnyTimes().Return()
			mhs.EXPECT().Run("9091").Return(nil)

			cmd = command.NewCommand(mhs)
		})

		Context("With FileStore args", func() {
			It("should run", func() {
				// remove USE_CLOUD_STORAGE if set
				os.Setenv("USE_CLOUD_STORAGE", "")
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

		// TODO: panic without having GOOGLE_APPLICATION_CREDENTIALS
		// Context("With CloudStorageStore args", func() {
		// 	It("should run", func() {
		// 		args := os.Args
		// 		os.Args = []string{"gae-fitbit-go",
		// 			"--port", "9091",
		// 			"--fb-client-id", "fb-client-id",
		// 			"--fb-client-secret", "fb-client-secret",
		// 			"--gcal-sleep-cal-id", "gcal-sleep-cal-id",
		// 			"--gcal-activity-cal-id", "gcal-activity-cal-id",
		// 			"--gcal-client-id", "gcal-client-id",
		// 			"--gcal-client-secret", "gcal-client-secret",
		// 			"--use-cloud-storage",
		// 			"--cloud-storage-bucket-name", "bucket-name",
		// 		}
		// 		err := cmd.Run()
		// 		Expect(err).To(BeNil())
		// 		os.Args = args
		// 	})
		// })
	})
})
