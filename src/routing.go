package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {

	/// PUBLIC ROUTES: ///

	// Add the token endpoint so users can get a token
	e.POST("/token", createToken)

	// Redirects users to the doc
	e.GET("/", redirectToDocumentation)

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
