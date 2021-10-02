package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type User struct {
	UserID   uint64 `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
}

type UserDTO struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

func createUser(c echo.Context) error {
	user := new(User)

	err := c.Bind(user)
	if err != nil {
		fmt.Println(" [ERROR] Bind parsing failed:\n ", err)
		return c.String(http.StatusBadRequest, "Invalid User format received.")
	}

	err = InsertUser(user)
	if err != nil {
		if err == EmailInUseError(user.Email) {
			return c.String(http.StatusConflict, err.Error())
		}
		return c.String(http.StatusInternalServerError, "Failed to create user.")
	}

	return c.String(http.StatusOK, "User successfully created.")
}

func readAllUsers(c echo.Context) error {
	users, err := GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retreve users.")
	}

	dtoArr := make([]*UserDTO, 0)
	for _, u := range users {
		dtoArr = append(dtoArr, &UserDTO{
			UserID: u.UserID,
			Email:  u.Email,
			Domain: u.Domain,
			Name:   u.Name,
		})
	}

	return c.JSON(http.StatusOK, dtoArr)
}

func readUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid User ID provided: "+c.Param("user_id"))
	}
	user, err := GetUserWithID(uint64(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retreve user.")
	}

	uDTO := UserDTO{
		UserID: user.UserID,
		Email:  user.Email,
		Domain: user.Domain,
		Name:   user.Name,
	}

	return c.JSON(http.StatusOK, uDTO)
}

func updateUser(c echo.Context) error {
	// Get new User data from body
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		fmt.Println(" [ERROR] Bind parsing failed:\n ", err)
		return c.String(http.StatusBadRequest, "Invalid User format received.")
	}

	err = UpdateUser(user)
	if err != nil {
		if err == EmailInUseError(user.Email) {
			return c.String(http.StatusConflict, err.Error())
		}
		return c.String(http.StatusInternalServerError, "Failed to update the user.")
	}

	return c.String(http.StatusOK, "Successfully updated the user.")
}

func deleteUser(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid User ID provided: "+c.Param("user_id"))
	}
	id := uint64(idSigned)

	err = DeleteUser(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete user.")
	}

	return c.String(http.StatusOK, "User successfully deleted.")
}
