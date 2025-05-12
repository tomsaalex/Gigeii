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
}

func MakeNewAvailabilityEditComponent() *ViewModel {
	daysOfTheWeek := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dateToday := time.Now().Format("2006-01-02")

	return &ViewModel{
		daysOfTheWeek: daysOfTheWeek,
		dateToday:     dateToday,
	}
}

func (m *ViewModel) Render(ctx context.Context, w io.Writer) error {
	return modal.MakeNewModal("Edit Availability", editAvailabilityCompBody(m), editAvailabilityCompFooter()).
		Render(ctx, w)
}
