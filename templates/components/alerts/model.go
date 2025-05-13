package custalerts

import (
	"context"
	"io"
)

type alertType string

const (
	alertDanger alertType = "alert-danger"
)

type Model struct {
	msgList   []string
	alertType alertType
}

func MakeMultiLineAlertDanger(msg []string) *Model {
	return &Model{
		msgList:   msg,
		alertType: alertDanger,
	}
}

func MakeAlertDanger(msg string) *Model {
	return &Model{
		msgList:   []string{msg},
		alertType: alertDanger,
	}
}

func (vm *Model) Render(ctx context.Context, w io.Writer) error {
	return component(vm).Render(ctx, w)
}
