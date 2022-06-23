package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logme = logrus.New()
)

func entryLog() {

	//logme.SetFormatter(&logrus.JSONFormatter{})
	logme.SetFormatter(&logrus.TextFormatter{})
	logme.SetOutput(os.Stdout)
	logme.SetLevel(logrus.InfoLevel)
}
