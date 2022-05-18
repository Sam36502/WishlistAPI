package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialization
	ConnectDB()
	InitRoutes(e)

	f, err := ioutil.ReadFile("/certs/fullchain.pem")
	if err != nil {
		fmt.Printf("Error: Failed to open file:\n  %s\n", err)
	}
	fmt.Println(string(f))

	e.Logger.Fatal(e.StartTLS(":"+os.Getenv("WISHLIST_API_PORTNUM"), os.Getenv("WISHLIST_SSL_CERT"), os.Getenv("WISHLIST_SSL_KEY")))
	//e.Logger.Fatal(e.Start(":2512"))
}
