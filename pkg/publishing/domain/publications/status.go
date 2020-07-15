package publications

import (
	"omics/pkg/shared/models"
	"time"
)

var (
	DRAFT            string = "draft"
	WAITING_APPROVAL        = "waiting-approval"
	APPROVED                = "approved"
	REJECTED                = "rejected"
)

type Status struct {
	code string
}

func (s Status) Is(code string) bool {
	return s.code == code
}

type ChangedBy struct {
	role   string
	userID models.ID
}

type StatusHistory struct {
	date      time.Time
	status    Status
	comment   string
	changedBy ChangedBy
}
