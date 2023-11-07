package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialization
	ConnectDB()
	InitRoutes(e)

	e.Logger.Fatal(e.Start(":" + os.Getenv("WISHLIST_API_PORTNUM")))
}
