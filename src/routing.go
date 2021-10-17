package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {

	/// PUBLIC ROUTES: ///

	// Add the token endpoint so users can get a token
	e.POST("/token", createToken)

	// Redirects users to the doc
	e.GET("/", redirectToDocumentation)

	// Displays the API's Status for checking
	e.GET("/status", statusPage)

	// User
	e.POST("/user", createUser)
	e.GET("/user", readAllUsers)
	e.GET("/user/:user_id", readUser)

	// Items
	e.GET("/user/:user_id/list", readAllItems)
	e.GET("/item/:item_id", readItem)

	/// PRIVATE ROUTES: ///
	/// Adds the TokenValidator middleware so it'll automatically check for tokens for every request
	/// All Routes defined below this call require a token to execute
	///// Currently not working. Just adding TokenValidator to all Private Routes
	///// Also for future reference: TokenValidator -> checks if ur logged in; AuthValidator -> checks if you are the user you're editing
	//e.Use(TokenValidator)

	// User
	e.PUT("/user/:user_id", updateUser, TokenValidator, AuthValidator)
	e.DELETE("/user/:user_id", deleteUser, TokenValidator, AuthValidator)

	// Item
	e.POST("/user/:user_id/list", createItem, TokenValidator, AuthValidator)
	e.PUT("/user/:user_id/list/:item_id", updateItem, TokenValidator, AuthValidator)
	e.DELETE("/user/:user_id/list/:item_id", deleteItem, TokenValidator, AuthValidator)

}

func redirectToDocumentation(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "https://www.pearcenet.ch/wishlist/doc.html")
}

func statusPage(c echo.Context) error {
	dbStatus := "<td class='red'>FAILED</td>"
	if IsDatabaseOnline() {
		dbStatus = "<td class='green'>CONNECTED</td>"
	}

	eaStatus := "<td class='red'>INACCESSIBLE</td>"
	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get("https://" + os.Getenv("WISHLIST_API_DOMAIN") + ":" + os.Getenv("WISHLIST_API_PORTNUM") + "/user")
	if err == nil && resp.StatusCode == 200 {
		eaStatus = "<td class='green'>ACCESSIBLE</td>"
		resp.Body.Close()
	}

	return c.HTML(
		http.StatusOK,
		"<html><head><title>Wishlist API</title>"+
			"<style>body {width: 50%;margin: auto;padding: 50px;} * {font-family: sans-serif;} .red {color: #FF4136;} .green {color: #2ECC40}</style>"+
			"</head><body><h1>Wishlist API</h1><h3>It's working!</h3><br><h3>Information:</h3><table>"+

			"<tr><td>Database Connection:</td>"+dbStatus+"</tr>"+
			"<tr><td>External Access:</td>"+eaStatus+"</tr>"+

			"</table></body></html>",
	)
}
