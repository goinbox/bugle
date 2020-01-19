#!/bin/bash

curDir=`dirname $0`
cd $curDir/../
prjHome=`pwd`

if [ ! -d $prjHome/logs ]
then
    mkdir -p $prjHome/logs
fi

cd $prjHome/main

go run demo.go --prj-home=$prjHome $@
