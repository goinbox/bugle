package command

import (
	"errors"
	"reflect"

	"github.com/goinbox/bugle/core"
	"github.com/goinbox/gohttp/router"
)

const (
	CmdNameRequest = "request"
)

func init() {
	register(CmdNameRequest, newRequestCommand)
}

var Router router.Router

func newRequestCommand() ICommand {
	rc := &RequestCommand{
		baseCommand: NewBaseCommand(),
		router:      Router,
	}

	rc.AddMustHaveArgs("env", "route-path").
		SetRunFunc(rc.run)

	rc.Fs.StringVar(&rc.routePath, "route-path", "", "path for route")
	rc.Fs.BoolVar(&rc.dryRun, "dry-run", false, "dry run")

	return rc
}

type RequestCommand struct {
	*baseCommand

	routePath string
	dryRun    bool

	router router.Router
}

func (c *RequestCommand) run() error {
	if c.router == nil {
		return errors.New("you must set router first")
	}

	route := c.router.FindRoute(c.routePath)
	if route == nil {
		return errors.New("no route for " + c.routePath)
	}

	action := route.NewActionFunc.Call(c.makeArgsValues())[0].Interface().(core.Action)

	action.Before()
	action.Run()
	action.After()
	action.Destruct()

	return nil
}

func (c *RequestCommand) makeArgsValues() []reflect.Value {
	params := &core.ActionParams{
		Env:       c.Env,
		RoutePath: c.routePath,
		DryRun:    c.dryRun,
		VarConf:   c.VarConf,
		Args:      c.ExtArgs,
	}

	return []reflect.Value{reflect.ValueOf(params)}
}
