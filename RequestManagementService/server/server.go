package server

import (
	requestHandler "RequestManagementService/handlers/request"
	"RequestManagementService/services/DatabaseService/database"
	"RequestManagementService/services/DatabaseService/datastore/request"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer() {
	e.Use(middleware.Logger())
	db, err := database.GetConnection()
	if err != nil {
		e.Logger.Errorf("Failed to connect to database : %v", err.Error())
	}

	reqDataStore := request.New(db)
	adsHandler := requestHandler.New(reqDataStore)
	requestsRoutes(e, adsHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
