package fitbit2gcal

type FitbitClient interface{
  GetSleepData() error
  GetActivityData() error
}
