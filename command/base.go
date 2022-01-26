package command

import (
	"errors"
	"flag"
	"strings"

	"github.com/goinbox/bugle/core"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
)

type ICommand interface {
	Run(args []string) error
}

type newCommandFunc func() ICommand

var commandTable = make(map[string]newCommandFunc)

func register(name string, ncf newCommandFunc) {
	commandTable[name] = ncf
}

func NewCommandByName(name string) ICommand {
	ncf, ok := commandTable[name]
	if !ok {
		return nil
	}

	return ncf()
}

type runFunc func() error

type baseCommand struct {
	Fs *flag.FlagSet

	Env     string
	ExtArgs map[string]string

	VarConf *core.VarConf

	mustHaveArgs map[string]bool
	existArgs    map[string]bool

	rf runFunc
}

func NewBaseCommand() *baseCommand {
	return &baseCommand{
		Fs:      new(flag.FlagSet),
		ExtArgs: make(map[string]string),

		mustHaveArgs: make(map[string]bool),
		existArgs:    make(map[string]bool),
	}
}

func (c *baseCommand) AddMustHaveArgs(names ...string) *baseCommand {
	for _, name := range names {
		c.mustHaveArgs[name] = true
	}

	return c
}

func (c *baseCommand) SetRunFunc(rf runFunc) *baseCommand {
	c.rf = rf

	return c
}

func (c *baseCommand) Run(args []string) error {
	err := c.parseArgs(args)
	if err != nil {
		return err
	}

	for name, _ := range c.mustHaveArgs {
		_, ok := c.existArgs[name]
		if !ok {
			return errors.New("Must have arg " + name)
		}
	}

	if c.Env != "" {
		err = c.parseVars()
		if err != nil {
			return err
		}
	}

	return c.rf()
}

func (c *baseCommand) parseArgs(args []string) error {
	c.Fs.StringVar(&c.Env, "env", "", "env name")

	err := c.Fs.Parse(args)
	if err != nil {
		return err
	}

	c.Fs.Visit(func(f *flag.Flag) {
		c.existArgs[f.Name] = true
	})

	for _, str := range c.Fs.Args() {
		item := strings.Split(str, "=")
		if len(item) == 2 {
			c.ExtArgs[item[0]] = item[1]
		}
	}

	return nil
}

func (c *baseCommand) parseVars() error {
	gvc, err := core.NewVarConf(core.GlobalVarPath(), c.ExtArgs)
	if err != nil {
		return err
	}

	err = gvc.Parse()
	if err != nil {
		return err
	}

	evc, err := core.NewVarConf(core.EnvVarPath(c.Env), c.ExtArgs)
	if err != nil {
		return err
	}

	err = evc.Parse()
	if err != nil {
		return err
	}

	path := core.TmpVarPath(c.Env)
	if gomisc.FileExist(path) {
		tvc, _ := core.NewVarConf(path, c.ExtArgs)
		err = tvc.Parse()
		if err == nil {
			c.VarConf = core.MergeVarConfs(gvc, evc, tvc)
			return nil

		}
		core.Logger.Error("tmp_var error", golog.ErrorField(err))
	}

	c.VarConf = core.MergeVarConfs(gvc, evc)

	return nil
}
