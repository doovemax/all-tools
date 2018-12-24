package main

import (
	"github.com/doovemax/all-tool/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		logrus.Errorln(err)
	}
}
