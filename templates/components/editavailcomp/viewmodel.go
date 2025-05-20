package editavailcomp

import (
	"context"
	"io"
	"time"

	"example.com/templates/components/modal"
)

type ViewModel struct {
	daysOfTheWeek []string
	dateToday     string

	modalVersion bool
}

func MakeNewAvailabilityEditComponent(modalVersion bool) *ViewModel {
	daysOfTheWeek := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dateToday := time.Now().Format("2006-01-02")

	return &ViewModel{
		daysOfTheWeek: daysOfTheWeek,
		dateToday:     dateToday,
		modalVersion:  modalVersion,
	}
}

func (m *ViewModel) Render(ctx context.Context, w io.Writer) error {
	if m.modalVersion {
		return modal.MakeNewModal("Edit Availability", EditAvailabilityCompBody(m), editAvailabilityCompFooter()).
			Render(ctx, w)
	} else {
		return EditAvailabilityCompBody(m).Render(ctx, w)
	}
}
