package core

import (
	"github.com/goinbox/gomisc"
	"github.com/goinbox/shell"

	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	FuncValueArg  = "ARG"
	FuncValueEnv  = "ENV"
	FuncValueMath = "MATH"
)

var VarDir string
var TmpVarDir string

func GlobalVarPath() string {
	return VarDir + "/global.json"
}

func EnvVarPath(env string) string {
	return VarDir + "/" + env + ".json"
}

func TmpVarPath(env string) string {
	return TmpVarDir + "/" + env + ".json"
}

var definedValueRegex = regexp.MustCompile("{{([^}]+)}}")
var funcValueRegex = regexp.MustCompile("__([A-Z]+)__\\(([^)]+)\\)")
var mathFuncRegex = regexp.MustCompile("([0-9]+)([+\\-*/])([0-9]+)")

type VarConf struct {
	Vars map[string]string

	confPath       string
	user           string
	extArgs        map[string]string
	unparsedValues map[string]string
}

func NewVarConf(confPath string, extArgs map[string]string) (*VarConf, error) {
	if !gomisc.FileExist(confPath) {
		return nil, errors.New("VarConf not exist: " + confPath)
	}

	return &VarConf{
		Vars: make(map[string]string),

		confPath:       confPath,
		user:           os.Getenv("USER"),
		extArgs:        extArgs,
		unparsedValues: make(map[string]string),
	}, nil
}

func MergeVarConfs(vcs ...*VarConf) *VarConf {
	cnt := len(vcs)
	if cnt == 0 {
		return nil
	}

	vc := vcs[0]
	for i := 1; i < cnt; i++ {
		for key, value := range vcs[i].Vars {
			vc.Vars[key] = value
		}
	}

	return vc
}

func (vc *VarConf) Parse() error {
	var varConfJson map[string]interface{}

	err := gomisc.ParseJsonFile(vc.confPath, &varConfJson)
	if err != nil {
		return err
	}

	for key, item := range varConfJson {
		vs, err := vc.parseVarJsonItemtoString(item)
		if err != nil {
			return err
		}

		vc.unparsedValues[key] = vs
	}

	for len(vc.unparsedValues) > 0 {
		for key, value := range vc.unparsedValues {
			value, delay, err := vc.ParseValueByDefined(value)
			if err != nil {
				return err
			}

			if !delay {
				vs, err := vc.parseValueByFunc(value)
				if err != nil {
					return err
				}

				vc.Vars[key] = vs
				delete(vc.unparsedValues, key)
			}
		}
	}

	return nil
}

/**
* return parsed value and whether delay parsed, if delay parsed, bool is true
 */
func (vc *VarConf) ParseValueByDefined(value string) (string, bool, error) {
	matches := definedValueRegex.FindAllStringSubmatch(value, -1)

	if len(matches) == 0 {
		return value, false, nil
	}

	var rs []string
	for _, item := range matches {
		k := item[1]
		vs, ok := vc.Vars[k]
		if ok {
			Logger.Debug("find var k:" + k + " in Vars")
		} else {
			_, ok := vc.unparsedValues[k]
			if ok {
				return "", true, nil
			}
			return "", false, errors.New("Undefined var: " + k)
		}

		rs = append(rs, item[0])
		rs = append(rs, vs)
	}

	return strings.NewReplacer(rs...).Replace(value), false, nil
}

func (vc *VarConf) parseVarJsonItemtoString(item interface{}) (string, error) {
	var r string

	switch item.(type) {
	case string:
		r = item.(string)
	case map[string]interface{}:
		mv := item.(map[string]interface{})
		v, ok := mv[vc.user]
		if !ok {
			v = mv["default"]
		}
		r = v.(string)
	default:
		return "", errors.New("item's type not support")
	}

	return strings.TrimSpace(r), nil
}

func (vc *VarConf) parseValueByFunc(value string) (string, error) {
	match := funcValueRegex.FindStringSubmatch(value)
	var err error

	if len(match) != 0 {
		switch match[1] {
		case FuncValueArg:
			Logger.Debug("parse value: " + value + " by arg func")
			value, err = vc.parseByArgFunc(match[2])
		case FuncValueEnv:
			Logger.Debug("parse value: " + value + " by Env func")
			value, err = vc.parseByEnvFunc(match[2])
		case FuncValueMath:
			Logger.Debug("parse value: " + value + " by math func")
			value, err = vc.parseByMathFunc(match[2])
		default:
			err = errors.New("Not support func " + match[1])
		}
	}

	return value, err
}

func (vc *VarConf) parseByArgFunc(argName string) (string, error) {
	vs, ok := vc.extArgs[argName]
	if !ok {
		return "", errors.New("Not has arg " + argName)
	}

	return vs, nil
}

func (vc *VarConf) parseByEnvFunc(envName string) (string, error) {
	vs, find := os.LookupEnv(envName)
	if find {
		return vs, nil
	}

	vs = string(shell.RunCmd("echo $" + envName).Output)
	vs = strings.TrimSpace(vs)
	if vs == "" {
		return "", errors.New("Not has Env " + envName)
	}

	return vs, nil
}

func (vc *VarConf) parseByMathFunc(express string) (string, error) {
	match := mathFuncRegex.FindStringSubmatch(express)

	if len(match) == 0 {
		return "", errors.New("Invalid match express " + express)
	}

	lv, _ := strconv.ParseInt(match[1], 10, 64)
	rv, _ := strconv.ParseInt(match[3], 10, 64)
	var value int64

	switch match[2] {
	case "+":
		value = lv + rv
	case "-":
		value = lv - rv
	case "*":
		value = lv * rv
	case "/":
		value = lv / rv
	}

	return strconv.FormatInt(value, 10), nil
}
