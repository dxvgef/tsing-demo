package action

import (
	"errors"

	"github.com/dxvgef/tsing"
)

type Index struct{}

func (Index) Demo(ctx tsing.Context) error {
	return ctx.Event(errors.New("test"))
}
