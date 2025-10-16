package application

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	managementuc "usp-management-device-api/business/management_uc"
	"usp-management-device-api/common/logging"
	httpcontroller "usp-management-device-api/controller/http"
	"usp-management-device-api/infras/enviroments"
	miniostore "usp-management-device-api/infras/minio"
	uspstore "usp-management-device-api/infras/sql_store"
	middleware "usp-management-device-api/middlewares"

	"github.com/gin-gonic/gin"
)

// Logrus level
// -- Trace
// -- Debug
// -- Info
// -- Warn
// -- Error
// -- Fatal
// -- Panic

func StartApplication(buildTime, version string) {
	// global enviroment store
	globalStore := enviroments.NewStore()
	globalStore.Load()

	fmt.Println("Application start from here...")
	globalStore.Print()

	// signal hook setup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	logging.SetLevel(globalStore.GetAppLogLevel())

	minioDB := minioConnection(
		globalStore.GetMinIOEndpoint(),
		globalStore.GetMinIOAccessKey(),
		globalStore.GetMinIOSecretKey())

	uspDB := mustConnectionSQL(
		globalStore.GetSQLHost(),
		globalStore.GetSQLUSPUser(),
		globalStore.GetSQLUSPPass(),
		globalStore.GetSQLUSPDB(),
		globalStore.GetSQLPort(),
		globalStore.GetAppName())

	go func() {
		<-c
		db, _ := uspDB.DB()
		db.Close()
		os.Exit(0)
	}()

	uspStore := uspstore.NewStore(uspDB)
	minioStore := miniostore.NewMinioStore(minioDB)

	mu := managementuc.NewManagementUsecase(uspStore, minioStore)

	router := gin.Default()

	if globalStore.GetAppDeployEnv() == "prod" {
		gin.SetMode(gin.ReleaseMode)
		router.Use(middleware.Recovery())
	}

	controller := httpcontroller.NewHTTPController(mu)
	controller.SetupRoute(router, buildTime, version)
	if err := router.Run(":" + strings.ReplaceAll(globalStore.GetAppPort(), ":", "")); err != nil {
		logging.Fatalf("Failed to start server: %v", err)
	}
}
