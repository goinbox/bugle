package core

import (
	"strconv"
	"strings"

	"github.com/goinbox/gohttp/httpclient"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
)

type Action interface {
	Name() string

	SetValue(key string, value interface{})
	Value(key string) interface{}

	Before()
	Run()
	After()
	Destruct()
}

type ActionParams struct {
	env       string
	routePath string
	dryRun    bool

	varConf *VarConf
	args    map[string]string
}

type BaseAction struct {
	params *ActionParams
	data   map[string]interface{}

	LogRequestBody  bool
	LogResponseBody bool

	savedVars map[string]string
}

func NewBaseAction(params *ActionParams) *BaseAction {
	return &BaseAction{
		params: params,
		data:   make(map[string]interface{}),
	}
}

func (a *BaseAction) SetValue(key string, value interface{}) {
	a.data[key] = value
}

func (a *BaseAction) Value(key string) interface{} {
	return a.data[key]
}

func (a *BaseAction) Env() string {
	return a.params.env
}

func (a *BaseAction) RoutePath() string {
	return a.params.routePath
}

func (a *BaseAction) DryRun() bool {
	return a.params.dryRun
}

func (a *BaseAction) VarValue(name string) string {
	return a.params.varConf.Vars[name]
}

func (a *BaseAction) ArgValue(name string) string {
	return a.params.args[name]
}

func (a *BaseAction) Before() {
}

func (a *BaseAction) After() {
	if len(a.savedVars) > 0 {
		a.saveVars()
	}
}

func (a *BaseAction) Destruct() {
}

func (a *BaseAction) DoRequest(req *httpclient.Request) (*httpclient.Response, error) {
	defer func() {
		Logger.Notice("=====================")
	}()

	Logger.Notice("start request")
	Logger.Info(req.Method + " " + req.Url)
	Logger.Info("host: " + req.Host)
	Logger.Info("request-header:")
	for field, vs := range req.Header {
		Logger.Info(field + ": " + strings.Join(vs, " "))
	}

	if a.LogRequestBody {
		Logger.Info("request-body:")
		Logger.Info(string(req.Body))
	}

	if a.params.dryRun {
		Logger.Notice("end with dry run")
		return nil, nil
	}

	client := httpclient.NewClient(httpclient.NewConfig(), Logger)
	resp, err := client.Do(req, 1)
	if err != nil {
		Logger.Error("end with error: " + err.Error())
		return nil, err
	}

	Logger.Notice("receive response")
	Logger.Info("time: " + resp.T.String())
	Logger.Info("status-code: " + strconv.Itoa(resp.StatusCode))
	Logger.Info("response-header:")
	for field, vs := range resp.Header {
		Logger.Info(field + ": " + strings.Join(vs, " "))
	}

	if a.LogResponseBody {
		Logger.Info("response-body:")
		Logger.Info(string(resp.Contents))
	}

	Logger.Notice("end request")

	return resp, nil
}

func (a *BaseAction) SaveVar(key, value string) {
	a.savedVars[key] = value
}

func (a *BaseAction) saveVars() {
	tm := make(map[string]string)
	path := TmpVarPath(a.params.env)

	if gomisc.FileExist(path) {
		err := gomisc.ParseJsonFile(path, &tm)
		if err != nil {
			Logger.Error("parse_tmp_var", golog.ErrorField(err))
		}
	}

	for k, v := range a.savedVars {
		tm[k] = v
	}

	err := gomisc.SaveJsonFile(path, tm)
	if err != nil {
		Logger.Error("save_tmp_var", golog.ErrorField(err))
	}
}
