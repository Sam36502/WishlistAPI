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

	e.Logger.Fatal(e.StartTLS(":"+os.Getenv("WISHLIST_API_PORTNUM"), os.Getenv("WISHLIST_SSL_CERT"), os.Getenv("WISHLIST_SSL_KEY")))
	//e.Logger.Fatal(e.Start(":2512"))
}
