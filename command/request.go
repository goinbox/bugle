package command

import (
	"errors"
	"github.com/goinbox/bugle/core"
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/router"
	"reflect"
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

func (rc *RequestCommand) run() error {
	if rc.router == nil {
		return errors.New("you must set router first")
	}

	route := rc.router.FindRoute(rc.routePath)
	if route == nil {
		return errors.New("no route for " + rc.routePath)
	}

	context := route.Cl.NewActionContext(nil, nil)
	context.SetValue(core.ContextValueKeyVarConf, rc.VarConf)
	context.SetValue(core.ContextValueKeyArgs, rc.ExtArgs)
	context.SetValue(core.ContextValueKeyEnv, rc.Env)
	context.SetValue(core.ContextValueKeyRoutePath, rc.routePath)
	context.SetValue(core.ContextValueKeyDryRun, rc.dryRun)

	context.BeforeAction()
	route.ActionValue.Call(rc.makeArgsValues(context, route.Args))
	context.AfterAction()

	return nil
}

func (rc *RequestCommand) makeArgsValues(context controller.ActionContext, args []string) []reflect.Value {
	argsValues := make([]reflect.Value, len(args)+1)
	argsValues[0] = reflect.ValueOf(context)
	for i, arg := range args {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}
