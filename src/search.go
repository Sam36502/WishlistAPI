package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.deanishe.net/fuzzy"
)

type UserList []*User

// Implementing fuzzy.Sortable to sort the user object
func (u UserList) Keywords(i int) string {
	return u[i].Name + " " + u[i].Email
}

// Default sort.Interface methods
func (u UserList) Len() int      { return len(u) }
func (u UserList) Swap(i, j int) { u[i], u[j] = u[j], u[i] }

// Less is used as a tie-breaker when fuzzy match score is the same.
func (u UserList) Less(i, j int) bool {
	return (u[i].Name + " " + u[i].Email) < (u[j].Name + " " + u[j].Email)
}

// Searches for a given name in the users
func SearchUsers(c echo.Context) error {

	// Get Search Query
	if !c.QueryParams().Has("search") {
		return c.String(http.StatusBadRequest, "Bad Request. 'search' query parameter required.")
	}
	query := c.QueryParam("search")

	// Get all users with the query in their names
	allUsers, err := GetUsersByNameOrEmail(query)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve users.")
	}

	// Sort results by how close they are to the search string
	// Using https://go.deanishe.net/fuzzy
	fuzzy.Sort(UserList(allUsers), query)

	return c.JSON(http.StatusOK, allUsers)
}
