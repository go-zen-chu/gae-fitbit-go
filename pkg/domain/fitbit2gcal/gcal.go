package fitbit2gcal

type GCalClient interface {
  PostSchedule() error
}
