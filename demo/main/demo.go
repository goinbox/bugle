package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goinbox/bugle"
	"github.com/goinbox/golog"
	"github.com/goinbox/router"

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

	w, _ := golog.NewFileWriter("/dev/stdout", 0)
	logger := golog.NewSimpleLogger(w, golog.NewSimpleFormater())

	bugle.SetLogger(logger)
	err := bugle.SetVarDir(prjHome+"/data/var", prjHome+"/tmp/var")
	if err != nil {
		fmt.Println("bugle.SetVarDir error:", err)
		os.Exit(1)
	}

	r := router.NewRouter()
	r.MapRouteItems(
		new(demo.Controller),
	)

	bugle.SetRouter(r)

	err = bugle.Run(os.Args[2:])
	if err != nil {
		logger.Error("bugle.Run error", golog.ErrorField(err))
	}
}
