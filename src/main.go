package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialization
	ConnectDB()
	loadClients()
	InitRoutes(e)

	// Get Port from Environment

	e.Logger.Fatal(e.Start(":" + os.Getenv("WISHLIST_API_PORTNUM")))
}
