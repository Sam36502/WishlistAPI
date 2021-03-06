package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Status struct {
	StatusID    uint64 `json:"status_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Link struct {
	LinkID    uint64 `json:"link_id"`
	Text      string `json:"text"`
	Hyperlink string `json:"hyperlink"`
}

type Item struct {
	ItemID         uint64   `json:"item_id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	ReservedByUser *UserDTO `json:"reserved_by,omitempty"`
	Status         Status   `json:"status"`
	Price          int      `json:"price"`
	User           UserDTO  `json:"user"` // Using the DTO, because Item doesn't need the password and shouldn't display it
	Links          []Link   `json:"links"`
}

func readAllItems(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid User ID provided: \""+c.Param("user_id")+"\"")
	}
	id := uint64(idSigned)

	items, err := GetAllItems(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retreve items.")
	}
	return c.JSON(http.StatusOK, items)
}

func createItem(c echo.Context) error {
	// Get new Item data from body
	item := new(Item)
	err := c.Bind(item)
	if err != nil {
		fmt.Println(" [ERROR] Bind parsing failed:\n ", err)
		return c.String(http.StatusBadRequest, "Invalid Item format received.")
	}

	// Using User ID from URL to save redundancy
	idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid User ID provided: \""+c.Param("user_id")+"\"")
	}
	item.User.UserID = uint64(idSigned)

	// Insert Item
	err = InsertItem(item)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to add item.")
	}

	return c.String(http.StatusOK, "Successfully added item")
}

func readItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Item ID provided: "+c.Param("item_id"))
	}
	item, err := GetItemWithID(uint64(id))
	if err != nil {
		return c.String(http.StatusNotFound, "Item not found.")
	}
	return c.JSON(http.StatusOK, item)
}

func updateItem(c echo.Context) error {
	// Get new Item data from body
	item := new(Item)
	item.Price = -1 // Set so that we can detect if the user set it to zero on purpose
	err := c.Bind(item)
	if err != nil {
		fmt.Println(" [ERROR] Bind parsing failed:\n ", err)
		return c.String(http.StatusBadRequest, "Invalid Item format received.")
	}

	// Using Item ID from URL to save redundancy
	idSigned, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Item ID provided: \""+c.Param("item_id")+"\"")
	}
	item.ItemID = uint64(idSigned)

	err = UpdateItem(item)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update the item.")
	}

	return c.String(http.StatusOK, "Successfully updated the item.")
}

func deleteItem(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Item ID provided: \""+c.Param("item_id")+"\"")
	}
	id := uint64(idSigned)

	err = DeleteItem(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete item.")
	}

	return c.String(http.StatusOK, "Item successfully deleted.")
}

func reserveItem(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Item ID provided: \""+c.Param("item_id")+"\"")
	}
	item_id := uint64(idSigned)

	// Get currently logged in token's ID
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve JWT data from middleware",
		}
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve user claims from middleware JWT",
		}
	}
	loggedInUser, err := GetUserWithEmail(claims.Email)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Token used is for a user that doesn't exist anymore",
		}
	}

	err = UpdateItem(&Item{
		ItemID: item_id,
		Status: Status{
			StatusID: 2,
		},
		ReservedByUser: &UserDTO{
			UserID: loggedInUser.UserID,
		},
		Price: -1, // Must be specifically set to not delete price
	})
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

func unreserveItem(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Item ID provided: \""+c.Param("item_id")+"\"")
	}
	item_id := uint64(idSigned)

	// Get currently logged in token's ID
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve JWT data from middleware",
		}
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve user claims from middleware JWT",
		}
	}
	loggedInUser, err := GetUserWithEmail(claims.Email)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Token used is for a user that doesn't exist anymore",
		}
	}

	item, err := GetItemWithID(item_id)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// Another user is trying to unreserve the item
	if item.ReservedByUser.UserID != loggedInUser.UserID {
		return c.NoContent(http.StatusForbidden)
	}

	err = UpdateItem(&Item{
		ItemID: item_id,
		Status: Status{
			StatusID: 1,
		},
		ReservedByUser: nil,
		Price:          -1, // Must be specifically set to not delete price
	})
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}
