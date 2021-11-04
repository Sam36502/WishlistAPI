package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Status struct {
	StatusID    uint64
	Name        string
	Description string
}

type Link struct {
	LinkID    uint64
	Text      string
	Hyperlink string
}

type Item struct {
	ItemID      uint64
	Name        string
	Description string
	Status      Status
	Price       int
	User        UserDTO // Using the DTO, because Item doesn't need the password and shouldn't display it
	Links       []Link
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
