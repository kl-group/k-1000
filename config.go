package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	Config = viper.New()
)

func entryConfig() {

	Config.SetConfigFile(flagOptions.Config)
	configDefault()
	err := Config.ReadInConfig()
	if err != nil {
		logme.Error(err.Error())
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
}

func configDefault() {
	Config.SetDefault("asterisk.extension.purge", false)
	Config.SetDefault("asterisk.extension.attribute", "name")
	Config.SetDefault("asterisk.extension.attribute-regexp", "^SIP-[0-9]{3}$")
	Config.SetDefault("asterisk.extension.owner.attribute", "managedBy")
	Config.SetDefault("asterisk.extension.owner.attributeType", "CN")

	Config.SetDefault("asterisk.ringgroup.purge", false)
	Config.SetDefault("asterisk.ringgroup.attribute", "name")
	Config.SetDefault("asterisk.ringgroup.attribute-regexp", "^SIP-GRP-[0-9]{4}$")

	Config.SetDefault("default.asterisk.extension.name", "no name extension")
	Config.SetDefault("default.asterisk.ringgroup.name", "no name ring group")
}
