package timeprovider

import "time"

type TimeProvider struct{}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (t *TimeProvider) Now() time.Time {
	return time.Now()
}

func (t *TimeProvider) NowWithAddDuration(duration time.Duration) time.Time {
	return t.Now().Add(duration)
}
