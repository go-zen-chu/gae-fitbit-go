//go:generate mockgen -source fitbit.go -destination mock_fitbit.go
package fitbit2gcal

type FitbitClient interface {
	GetSleepData(dateStr string) (*Sleep, error)
	GetActivityData(dateStr string) (*Activity, error)
}
