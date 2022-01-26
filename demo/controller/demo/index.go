package demo

import (
	"net/http"

	"github.com/goinbox/gohttp/httpclient"

	"github.com/goinbox/bugle/core"
)

type indexAction struct {
	*core.BaseAction
}

func newIndexAction(params *core.ActionParams) *indexAction {
	return &indexAction{core.NewBaseAction(params)}
}

func (a *indexAction) Name() string {
	return "index"
}

func (a *indexAction) Run() {
	req, _ := httpclient.NewRequest(http.MethodGet, "http://127.0.0.1:8010", []byte("hello"), "", map[string]string{"Dev-By": "ligang"})

	resp, err := a.DoRequest(req)
	if a.DryRun() {
		return
	}

	if err == nil && resp.StatusCode == http.StatusOK {
		for k, v := range a.Args() {
			a.SaveVar(k, v)
		}
	}
}
