package action

import (
	"errors"

	"github.com/dxvgef/tsing"
)

type Manager struct{}

func (Manager) Auth(ctx tsing.Context) error {

	// panic("haha")
	return ctx.Event(errors.New("test"))
}
