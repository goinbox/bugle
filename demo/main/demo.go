package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goinbox/bugle"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/golog"

	"github.com/goinbox/bugle/demo/controller/demo"
)

func main() {
	var prjHome string

	flag.StringVar(&prjHome, "prj-home", "", "prj-home absolute path")
	flag.Parse()

	prjHome = strings.TrimRight(prjHome, "/")
	if prjHome == "" {
		fmt.Println("Missing flag prj-home: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger := golog.NewSimpleLogger(
		golog.NewConsoleWriter(),
		golog.NewConsoleFormater(new(golog.NoopFormater)))

	w, err := golog.NewFileWriter(prjHome+"/logs/request.log", 0)
	if err != nil {
		logger.Error([]byte(err.Error()))
		os.Exit(1)
	}
	rlogger := golog.NewSimpleLogger(w, golog.NewSimpleFormater())

	bugle.SetLogger(logger, rlogger)
	err = bugle.SetVarDir(prjHome+"/data/var", prjHome+"/tmp/var")
	if err != nil {
		logger.Error([]byte(err.Error()))
	}

	r := router.NewSimpleRouter()
	r.MapRouteItems(
		new(demo.DemoController),
	)

	bugle.SetRouter(r)

	err = bugle.Run(os.Args[2:])
	if err != nil {
		logger.Error([]byte(err.Error()))
	}
}
