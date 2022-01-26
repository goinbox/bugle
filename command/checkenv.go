package command

import (
	"fmt"

	"github.com/goinbox/bugle/core"
)

const (
	CmdNameCheckEnv = "checkenv"
)

func init() {
	register(CmdNameCheckEnv, newCheckEnvCommand)
}

func newCheckEnvCommand() ICommand {
	rc := &CheckEnvCommand{
		baseCommand: NewBaseCommand(),
	}

	rc.AddMustHaveArgs("env").
		SetRunFunc(rc.run)

	return rc
}

type CheckEnvCommand struct {
	*baseCommand
}

func (rc *CheckEnvCommand) run() error {
	core.Logger.Warning("checkenv list")
	for name, value := range rc.VarConf.Vars {
		fmt.Println(name, value)
	}

	return nil
}
