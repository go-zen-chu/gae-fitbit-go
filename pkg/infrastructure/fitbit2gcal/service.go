package fitbit2gcal

import (
	"net/http"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
)

type service struct{}

func NewService() df2g.Service {
	return &service{}
}

func (s *service) HandleFitbit2GCal(w http.ResponseWriter, r *http.Request) {

}
