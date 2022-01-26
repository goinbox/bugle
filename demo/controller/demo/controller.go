package demo

import (
	"github.com/goinbox/bugle/core"
)

type Controller struct {
}

func (c *Controller) Name() string {
	return "demo"
}

func (c *Controller) IndexAction(params *core.ActionParams) *indexAction {
	return newIndexAction(params)
}
