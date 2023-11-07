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
	Name     string `json:"name"`
}

type UserDTO struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func createUser(c echo.Context) error {
	user := new(User)

	err := c.Bind(user)
	if err != nil {
		fmt.Println(" [ERROR] Bind parsing failed:\n ", err)
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "invalid_user_format",
			Message: "Invalid User format received.",
		})
	}

	err = InsertUser(user)
	if err != nil {
		if err == EmailInUseError(user.Email) {
			return c.JSON(http.StatusConflict, ErrorDTO{
				Code:    "email_in_use",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "user_create_fail",
			Message: "Failed to create user.",
		})
	}

	return c.JSON(http.StatusOK, ErrorDTO{
		Code:    "user_create_succ",
		Message: "User successfully created.",
	})
}

func readAllUsers(c echo.Context) error {
	users, err := GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "retrieve_all_user_fail",
			Message: "Failed to retreve users.",
		})
	}

	dtoArr := make([]*UserDTO, 0)
	for _, u := range users {
		dtoArr = append(dtoArr, &UserDTO{
			UserID: u.UserID,
			Email:  u.Email,
			Name:   u.Name,
		})
	}

	return c.JSON(http.StatusOK, dtoArr)
}

func readUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "invalid_user_id",
			Message: "Invalid User ID provided: " + c.Param("user_id"),
		})
	}
	user, err := GetUserWithID(uint64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "retrieve_user_fail",
			Message: "Failed to retreve user.",
		})
	}

	uDTO := UserDTO{
		UserID: user.UserID,
		Email:  user.Email,
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
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "invalid_user_format",
			Message: "Invalid User format received.",
		})
	}

	// Get ID from URL
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "invalid_user_id",
			Message: "Invalid User ID provided: " + c.Param("user_id"),
		})
	}
	user.UserID = uint64(id)

	err = UpdateUser(user)
	if err != nil {
		if err == EmailInUseError(user.Email) {
			return c.JSON(http.StatusConflict, ErrorDTO{
				Code:    "email_in_use",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "edit_user_fail",
			Message: "Failed to update the user.",
		})
	}

	return c.JSON(http.StatusOK, ErrorDTO{
		Code:    "edit_user_succ",
		Message: "Successfully updated the user.",
	})
}

func deleteUser(c echo.Context) error {
	idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "invalid_user_id",
			Message: "Invalid User ID provided: " + c.Param("user_id"),
		})
	}
	id := uint64(idSigned)

	err = DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "delete_user_fail",
			Message: "Failed to delete user.",
		})
	}

	return c.JSON(http.StatusOK, ErrorDTO{
		Code:    "delete_user_succ",
		Message: "User successfully deleted.",
	})
}

func userByEmail(c echo.Context) error {
	if !c.QueryParams().Has("email") {
		return c.JSON(http.StatusBadRequest, ErrorDTO{
			Code:    "param_required_email",
			Message: "Required query parameter 'email' is missing.",
		})
	}

	email := c.QueryParam("email")

	user, err := GetUserWithEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorDTO{
			Code:    "retrieve_user_fail",
			Message: "Failed to retrieve user.",
		})
	}

	userDTO := UserDTO{
		UserID: user.UserID,
		Email:  user.Email,
		Name:   user.Name,
	}

	return c.JSON(http.StatusOK, userDTO)
}
