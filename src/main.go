package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialization
	ConnectDB()
	loadClients()
	InitRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
