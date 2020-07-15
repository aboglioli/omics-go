package publications

import (
	"time"

	"omics/pkg/shared/models"
)

// Publication is an aggregate root
type Publication struct {
	id            models.ID
	name          Name
	synopsis      Synopsis
	authorID      models.ID
	statusHistory []StatusHistory
	pages         []Page
	statistics    Statistics
	contract      bool

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewPublication(
	id models.ID,
	authorID models.ID,
	name string,
	synopsis string,
	author Author,
) (*Publication, error) {
	errs := ErrValidation

	p := &Publication{
		id:       id,
		authorID: authorID,
		statusHistory: []StatusHistory{
			NewStatusHistory(DRAFT, ""),
		},
	}

	if err := p.SetName(name); err != nil {
		errs = errs.Merge(err)
	}

	if err := p.SetSynopsis(synopsis); err != nil {
		errs = errs.Merge(err)
	}

	if errs.ContextLen() > 0 {
		return nil, errs
	}

	return p, nil
}

func (p *Publication) Author() Author {
	return p.author
}

func (p *Publication) SetName(name string) error {
	if !p.canBeModified() {
		return Err
	}

	n, err := NewName(name)
	if err != nil {
		return err
	}

	p.name = n
}

func (p *Publication) SetSynopsis(synopsis string) error {
	if !p.canBeModified() {
		return Err
	}

	s, err := NewSynopsis(synopsis)
	if err != nil {
		return err
	}

	p.synopsis = s

	return nil
}

func (p *Publication) CurrentStatus() Status {
	return p.statusHistory[0].status // TODO
}

func (p *Publication) Approve(user User, comment string) error {
	if p.CurrentStatus().Is(WAITING_APPROVAL) {
		p.setStatus(user, APPROVED, comment)
		return nil
	}
	return Err
}

func (p *Publication) Reject(user User, comment string) error {
	if p.CurrentStatus().Is(WAITING_APPROVAL) {
		p.setStatus(user, REJECTED, comment)
		return nil
	}
	return Err
}

func (p *Publication) Publish() error {
	if p.CurrentStatus().Is(DRAFT) {
		p.setStatus(p.author.toUser(), WAITING_APPROVAL, "")
		return nil
	}
	return Err
}

func (p *Publication) MakeAsDraft() error {
	p.setStatus(p.author.toUser(), DRAFT, "")
	return Err
}

func (p *Publication) canBeModified() bool {
	return p.CurrentStatus().Is(DRAFT)
}

func (p *Publication) setStatus(user User, code, comment string) {
	p.statusHistory = append(p.statusHistory, StatusHistory{
		date: time.Now(),
		status: Status{
			code: code,
		},
		comment: comment,
		changedBy: ChangedBy{
			role:   user.Role,
			userID: user.ID,
		},
	})
}
