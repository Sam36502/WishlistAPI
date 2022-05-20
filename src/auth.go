package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	TOKEN_EXPIRY    = 30 * 24 * 60 * 60 // 30 Days
	ENV_SIGNING_KEY = "WISHLIST_TOK_SIGNING_KEY"
)

type AuthRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type TokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Handler to receive authentication requests and return tokens
func tokenHandler(c echo.Context) error {

	// Parse Auth request
	var authreq AuthRequestDTO
	err := c.Bind(&authreq)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid auth request format.",
		}
	}

	// Find User details
	usr, err := GetUserWithEmail(authreq.Email)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "No user with that e-mail address found.",
		}
	}

	// Hash request password
	rows, err := Database.Query("SELECT SHA1(UNHEX(SHA1(?)))", authreq.Password)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to Hash Password: ", err)
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to hash password while authenticating. Aborted.",
		}
	}
	defer rows.Close()
	hashedPassword := ""
	rows.Next()
	err = rows.Scan(&hashedPassword)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to retrieve hashed password: ", err)
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to hash password while authenticating. Aborted.",
		}
	}

	// Check auth details are valid
	if usr.Password != hashedPassword {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials provided",
		}
	}

	// Hash Token
	expiry := time.Now().Unix() + TOKEN_EXPIRY
	tok := generateToken(authreq.Email, expiry)
	signd, err := tok.SignedString([]byte(os.Getenv(ENV_SIGNING_KEY)))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to create token:\n  %s\n", err),
		}
	}

	return c.JSON(http.StatusOK, AuthResponseDTO{
		Token:     signd,
		ExpiresAt: time.Unix(expiry, 0).Local().Format(time.RFC3339),
	})
}

// Generates a JWT with the user's email and expiry time
func generateToken(email string, expiresAt int64) *jwt.Token {
	cl := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Email: email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
}

// Middleware to check if a token is valid
var TokenValidator = middleware.JWTWithConfig(middleware.JWTConfig{
	Claims:     &TokenClaims{},
	SigningKey: []byte(os.Getenv(ENV_SIGNING_KEY)),
})

// AuthValidator Middleware makes sure the user making the request is the user being altered
func AuthValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

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

		idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid User ID provided: '%s'", c.Param("user_id")),
			}
		}
		id := uint64(idSigned)

		// if user is not updating themselves, abort
		if id != loggedInUser.UserID {
			return &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "You are forbidden from changing/deleting other users than yourself",
			}
		}

		return next(c)
	}
}
