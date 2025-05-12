package modal

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

type Model struct {
	title       string
	modalBody   templ.Component
	modalFooter templ.Component
}

func MakeNewModal(title string, modalBody, modalFooter templ.Component) *Model {
	return &Model{
		title:       title,
		modalBody:   modalBody,
		modalFooter: modalFooter,
	}
}

func (m *Model) Render(ctx context.Context, w io.Writer) error {
	return component(m).Render(ctx, w)
}
