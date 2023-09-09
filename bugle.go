package bugle

import (
	"errors"
	"os"
	"strings"

	"github.com/goinbox/bugle/command"
	"github.com/goinbox/bugle/core"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
	"github.com/goinbox/router"
)

func SetLogger(logger golog.Logger) {
	core.Logger = logger
}

func SetVarDir(varDir string, tmpVarDir string) error {
	varDir = strings.TrimRight(varDir, "/")
	if !gomisc.DirExist(varDir) {
		return errors.New("var_dir not exists")
	}

	tmpVarDir = strings.TrimRight(tmpVarDir, "/")
	if !gomisc.DirExist(tmpVarDir) {
		err := os.MkdirAll(tmpVarDir, 0755)
		if err != nil {
			return err
		}
	}

	core.VarDir = varDir
	core.TmpVarDir = tmpVarDir

	return nil
}

func SetRouter(router router.Router) {
	command.Router = router
}

func Run(argsStartWithCmd []string) error {
	if len(argsStartWithCmd) == 0 {
		return errors.New("must start with arg command")
	}

	cmd := command.NewCommandByName(argsStartWithCmd[0])
	if cmd == nil {
		return errors.New("command not exists")
	}

	return cmd.Run(argsStartWithCmd[1:])
}
