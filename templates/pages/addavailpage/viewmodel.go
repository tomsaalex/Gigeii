package addavailpage

import (
	"context"
	"io"
	"time"
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
	return availabilityAddPage(vm).Render(ctx, w)
}
