package main

/*
 *	Oauth Code
 *
 *	Most of this is stolen from an old work project
 *	I am unsure of how secure it is, but it's far
 *	better than whatever wacky system I'd come up with.
 *
 *	It implements OpenID Connect client_credentials flow
 *	There's also likely a better way of implementing
 *	this in Go, but this is the fastest way I know of
 *
 *	I will now try to comment this code as much as I
 *	can remember to make it somewhat debuggable...
 */

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

// Manages and stores tokens
var manager *manage.Manager

// Manages and stores clients
var clientStore *store.ClientStore

// Handles communication
var oauth2server *server.Server

// Type stored in clientStore. implements oauth2.ClientInfo
type clientStruct struct {
	ID     string
	Secret string
	Domain string
}

func (c clientStruct) GetID() string     { return c.ID }
func (c clientStruct) GetSecret() string { return c.Secret }
func (c clientStruct) GetDomain() string { return c.Domain }
func (c clientStruct) GetUserID() string { return "" }

// Loads all the clients out of the database
func loadClients() {
	manager = manage.NewManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(accessGen{})

	clientStore = store.NewClientStore()
	manager.MapClientStorage(clientStore)

	// Load user info into clientStore
	fmt.Println(" [INFO] Loading Clients...")
	allUsers, err := GetAllUsers()
	if err != nil {
		fmt.Println(" [ERROR] Failed to load clients.")
	}
	for _, user := range allUsers {
		clientStore.Set(user.Email, clientStruct{
			ID:     user.Email,
			Secret: user.Password,
			Domain: user.Domain,
		})
	}
	fmt.Println(" [INFO] Finished Loading Clients.")

	// Configure what type of token to issue
	config := &server.Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{},
		AllowedGrantTypes:    []oauth2.GrantType{oauth2.ClientCredentials},
	}

	oauth2server = server.NewServer(config, manager)
	oauth2server.SetClientInfoHandler(clientInfoHandler)
	oauth2server.SetClientAuthorizedHandler(clientAuthHandler)
}

// Replaces the default handler, because it wouldn't hash the passwords
// before checking them
func clientInfoHandler(r *http.Request) (clientID, clientSecret string, err error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return "", "", errors.ErrInvalidClient
	}

	// Hash Password
	rows, err := Database.Query("SELECT PASSWORD(?) AS `password`", password)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to Hash Password.")
		return "", "", errors.ErrInvalidClient
	}
	hashedPassword := ""
	rows.Next()
	err = rows.Scan(&hashedPassword)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to retrieve hashed password.")
		return "", "", errors.ErrInvalidClient
	}

	return username, hashedPassword, nil
}

// Sorry, don't entirely remember what this is for, but the library needs it
// I think it makes sure the library knows what type of grant to lok out for
func clientAuthHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	_, err = manager.GetClient(clientID)
	if err != nil {
		return false, err
	}

	if grant != "client_credentials" {
		return false, err
	}

	return true, nil
}

type accessGen struct{}

// Generates a random 32-char Base64 string to act as the token
func (accessGen) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	// Generate a random 32-Char Base 64 string
	str := ""
	for i := 0; i < 32; i++ {
		str += strconv.FormatInt(int64(rand.Intn(15-1)), 16)
	}
	tkn := base64.StdEncoding.EncodeToString([]byte(str))

	return tkn, "", nil
}

// TokenValidator Middleware to check the provided token
func TokenValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := c.QueryParam("token")
		info, err := manager.LoadAccessToken(access)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Failed to authorise; invalid token provided.")
		}

		c.Set("client_id", info.GetClientID())
		return next(c)
	}
}

// AuthValidator Middleware makes sure the user making the request is the user being altered
func AuthValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Get currently logged in client's ID
		clientEmail := fmt.Sprint(c.Get("client_id"))
		loggedInUser, err := GetUserWithEmail(clientEmail)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Invalid User logged in: "+clientEmail)
		}

		idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid User ID provided: "+c.Param("user_id"))
		}
		id := uint64(idSigned)

		// if user is not updating themselves, abort
		if id != loggedInUser.UserID {
			return c.String(http.StatusForbidden, "You are forbidden from changing/deleting other users than yourself.")
		}

		return next(c)
	}
}

// Endpoint for getting new tokens
func createToken(c echo.Context) error {

	// Get User's email and password from Auth header
	id, password, ok := c.Request().BasicAuth()
	if !ok {
		// Failed to get auth header
		return c.String(http.StatusUnauthorized, "Failed to create token; no Authorization header provided.")
	}

	// Hash Password
	rows, err := Database.Query("SELECT PASSWORD(?) AS `password`", password)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to Hash Password.")
		return c.String(http.StatusUnauthorized, "Failed to hash password. Aborted.")
	}
	hashedPassword := ""
	rows.Next()
	err = rows.Scan(&hashedPassword)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed. Failed to retrieve hashed password.")
		return c.String(http.StatusUnauthorized, "Failed to hash password. Aborted.")
	}

	// Get Client from list of stored clients
	client, err := clientStore.GetByID(id)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to create token; no such user exists.")
	}

	// Check the stored Client's password matches that provided
	if client.GetSecret() != hashedPassword {
		return c.String(http.StatusUnauthorized, "Failed to create token; invalid credentials provided.")
	}

	return oauth2server.HandleTokenRequest(c.Response(), c.Request())
}
