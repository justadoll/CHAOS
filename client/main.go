package main

import (
	"github.com/justadoll/CHAOS/client/app"
	"github.com/justadoll/CHAOS/client/app/shared/environment"
	"github.com/justadoll/CHAOS/client/app/ui"
	"github.com/justadoll/CHAOS/client/app/utilities/network"
)

var (
	Version       = "dev"
	ServerPort    = ""
	ServerAddress = ""
	Token         = ""
)

func main() {
	ui.ShowMenu(Version, ServerAddress, ServerPort)

	app.NewApp(network.NewHttpClient(10),
		environment.Load(ServerAddress, ServerPort, Token)).Run()
}
