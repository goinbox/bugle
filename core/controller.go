package core

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/httpclient"
	"github.com/goinbox/gomisc"
)

const (
	ContextValueKeyVarConf   = "var_conf"
	ContextValueKeyArgs      = "args"
	ContextValueKeyEnv       = "env"
	ContextValueKeyRoutePath = "route_path"
	ContextValueKeyDryRun    = "dry_run"
)

type BaseContext struct {
	*controller.BaseContext

	VarConf *VarConf
	Args    map[string]string

	LogRequestBody  bool
	LogResponseBody bool

	env       string
	routePath string
	dryRun    bool
	savedVars map[string]string
}

func (bc *BaseContext) BeforeAction() {
	bc.VarConf = bc.Value(ContextValueKeyVarConf).(*VarConf)
	bc.Args = bc.Value(ContextValueKeyArgs).(map[string]string)
	bc.env = bc.Value(ContextValueKeyEnv).(string)
	bc.routePath = bc.Value(ContextValueKeyRoutePath).(string)
	bc.dryRun = bc.Value(ContextValueKeyDryRun).(bool)
}

func (bc *BaseContext) AfterAction() {
	if len(bc.savedVars) > 0 {
		bc.saveVars()
	}
}

func (bc *BaseContext) Env() string {
	return bc.env
}

func (bc *BaseContext) RoutePath() string {
	return bc.routePath
}

func (bc *BaseContext) DryRun() bool {
	return bc.dryRun
}

func (bc *BaseContext) DoRequest(req *httpclient.Request) (*httpclient.Response, error) {
	defer func() {
		RequestLogger.Notice([]byte("====================="))
	}()

	RequestLogger.Notice([]byte("start request"))
	RequestLogger.Info([]byte(req.Method + " " + req.Url))
	RequestLogger.Info([]byte("host: " + req.Host))
	RequestLogger.Info([]byte("request-header:"))
	for field, vs := range req.Header {
		RequestLogger.Info([]byte(field + ": " + strings.Join(vs, " ")))
	}

	if bc.LogRequestBody {
		RequestLogger.Info([]byte("request-body:"))
		RequestLogger.Info(req.Body)
	}

	if bc.dryRun {
		RequestLogger.Notice([]byte("end with dry run"))
		return nil, nil
	}

	client := httpclient.NewClient(httpclient.NewConfig(), Logger)
	resp, err := client.Do(req, 1)
	if err != nil {
		RequestLogger.Error([]byte("end with error: " + err.Error()))
		return nil, err
	}

	RequestLogger.Notice([]byte("receive response"))
	RequestLogger.Info([]byte("time: " + resp.T.String()))
	RequestLogger.Info([]byte("status-code: " + strconv.Itoa(resp.StatusCode)))
	RequestLogger.Info([]byte("response-header:"))
	for field, vs := range resp.Header {
		RequestLogger.Info([]byte(field + ": " + strings.Join(vs, " ")))
	}

	if bc.LogResponseBody {
		RequestLogger.Info([]byte("response-body:"))
		RequestLogger.Info(resp.Contents)
	}

	RequestLogger.Notice([]byte("end request"))

	return resp, nil
}

func (bc *BaseContext) SaveVar(key, value string) {
	bc.savedVars[key] = value
}

func (bc *BaseContext) saveVars() {
	tm := make(map[string]string)
	path := TmpVarPath(bc.env)

	if gomisc.FileExist(path) {
		err := gomisc.ParseJsonFile(path, &tm)
		if err != nil {
			ErrorLog("parse_tmp_var", err.Error())
		}
	}

	for k, v := range bc.savedVars {
		tm[k] = v
	}

	err := gomisc.SaveJsonFile(path, tm)
	if err != nil {
		ErrorLog("save_tmp_var", err.Error())
	}
}

type BaseController struct {
}

func (bc *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &BaseContext{
		BaseContext: controller.NewBaseContext(req, respWriter),

		LogRequestBody:  true,
		LogResponseBody: true,

		savedVars: make(map[string]string),
	}
}
