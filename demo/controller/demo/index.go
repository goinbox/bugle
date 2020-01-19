package demo

import (
	"net/http"

	"github.com/goinbox/gohttp/httpclient"

	"github.com/goinbox/bugle/core"
)

func (d *DemoController) IndexAction(context *core.BaseContext) {
	req, _ := httpclient.NewRequest(http.MethodGet, "http://127.0.0.1:8010", []byte("hello"), "", map[string]string{"Dev-By": "ligang"})

	resp, err := context.DoRequest(req)
	if context.DryRun() {
		return
	}

	if err == nil && resp.StatusCode == http.StatusOK {
		for k, v := range context.Args {
			context.SaveVar(k, v)
		}
	}
}
