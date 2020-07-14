package publications

import "time"

type Statistics struct {
	Likes       int
	Views       int
	UniqueViews int
	Rate        float64
	UpdatedAt   time.Time
}
