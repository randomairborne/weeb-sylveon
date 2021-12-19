package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api"
	"github.com/randomairborne/eevee/api/logger"
	"github.com/randomairborne/eevee/modules/coremodules"
	"strings"
)

var loadedModules = make(map[string]api.Module, 0)

func LoadModule(ds *discordgo.Session, modules []string) {
	loadedModules["core"] = &coremodules.Module{}
	for _, v := range modules {
		logger.Out().Printf("Loading %s\n", v)
		switch strings.ToLower(v) {
		// case "hugs":
		// 	loadedModules["coremodules"] = &coremodules.Module{}
		default:
			logger.Err().Printf("No module with name %s\n", v)
		}
	}
	println(loadedModules)
	for k, v := range loadedModules {
		v.Load(ds)
		logger.Out().Printf("Loaded %s\n", k)
	}
}
