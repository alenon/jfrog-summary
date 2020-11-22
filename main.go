package main

import (
	"github.com/alenon/jfrog-summary/commands"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "rt-summary"
	app.Description = "Artifactory summary live visualisation"
	app.Version = "v0.1"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetStorageCommand()}
}
