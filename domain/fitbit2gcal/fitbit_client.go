//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbit2gcal

type FitbitClient interface {
	GetSleepData(dateStr string) (*Sleep, error)
	GetActivityData(dateStr string) (*Activity, error)
}
