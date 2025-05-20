package addavailpage

import (
	"context"
	"io"
	"time"

	"example.com/templates/components/editavailcomp"
)

type ViewModel struct {
	daysOfTheWeek []string
	dateToday     string
}

func MakeNewAvailabilityAddPage() *ViewModel {
	daysOfTheWeek := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dateToday := time.Now().Format("2006-01-02")

	return &ViewModel{
		daysOfTheWeek: daysOfTheWeek,
		dateToday:     dateToday,
	}
}

func (vm *ViewModel) Render(ctx context.Context, w io.Writer) error {
	return availabilityAddPage(editavailcomp.MakeNewAvailabilityEditComponent(false)).Render(ctx, w)
}
