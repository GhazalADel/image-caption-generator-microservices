package server

import (
	"RequestManagementService/handlers/request"
	"github.com/labstack/echo/v4"
)

func requestsRoutes(e *echo.Echo, handler *request.RequestHandler) {
	e.POST("/requests", handler.AddRequest)
	e.GET("/requests/:id", handler.GetRequest)
}
