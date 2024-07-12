package main

import (
	"github.com/MasoudHeydari/eps-api/cmd"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
)

func main() {
	defer recoverPanic()
	unix.Umask(0002)
	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Info(err)
		os.Exit(1)
	}
}

func recoverPanic() {
	// Add comment for recoverPanic function
	if r := recover(); r != nil {
		logrus.Fatalf("recover.Error: %v\n", r)
	}
}
