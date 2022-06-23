package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type Opts struct {
	Config string `long:"config" short:"c" description:"Path to config file" default:"./config.yaml"`
	Daemon bool   `long:"daemon" short:"d" description:"Run as Daemon" `
}

var flagOptions Opts

var parser = flags.NewParser(&flagOptions, flags.Default)

func entryFlags() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Println(err.Error())
			os.Exit(1)
		}
	}

}
