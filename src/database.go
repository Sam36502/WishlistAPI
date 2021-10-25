/***********************************************************

		DATABASE FUNCTIONS

		It's pretty monolithic, because I don't like
		DB integration code and wanted it to all be
		hidden away in the corner.

		I tried to keep it somewhat maintainable by
		leaving copious navigation comments, but it'd
		be better if properly cleaned up.

		If I revisit this and want to clean up the
		code, I'd suggest putting this in its own
		module and grouping the functions by object
		(item, user, etc.)

		-2021

************************************************************/

package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB
var timerKillChannel = make(chan bool)

var (
	USERNAME = os.Getenv("WISHLIST_DB_USERNAME")
	PASSWORD = os.Getenv("WISHLIST_DB_PASSWORD")
	DATABASE = os.Getenv("WISHLIST_DB_DATABASE")
	HOSTNAME = os.Getenv("WISHLIST_DB_HOSTNAME")
)

// Tries to connect to the DB for 5 minutes before giving up
const CONN_TIMEOUT = 5 * time.Minute

func ConnectDB() {
	dataSourceStr := USERNAME + ":" + PASSWORD + "@tcp(" + HOSTNAME + ":3306)/" + DATABASE

	fmt.Printf(" [INFO] Connecting to Database...\n")
	var err error
	go DBConnTimer(timerKillChannel)
	for {
		Database, err = sql.Open("mysql", dataSourceStr)
		err = Database.Ping()
		if err == nil {
			timerKillChannel <- true
			break
		}
	}

	fmt.Printf(" [INFO] Successfully connected to Database!\n")
	Database.Exec("SET NAMES 'utf8mb4'")
}

// Pings database and returns true if it's online and can be connected to; otherwise false
func IsDatabaseOnline() bool {
	err := Database.Ping()
	return err == nil
}

func DBConnTimer(killChan chan bool) {
	kill := false
	for i := 0; i < 50; i++ {
		time.Sleep(CONN_TIMEOUT / 50)
		select {
		case kill = <-killChan:
		default:
		}
		if kill {
			return
		}
	}

	// Time ran out, quit program having failed to connect
	fmt.Printf(" [ERROR] Failed to connect to the database. Please check the connection and try again.\n")
	os.Exit(1)
}

/// USER FUNCTIONS: ///
//	GetAllUsers() -> []User, err
//	GetUserWithID( ID ) -> User, err
//	GetUserWithEmail( Email ) -> User, err
//	InsertUser( User ) -> err
//	UpdateUser( User ) -> err
//	DeleteUser( ID ) -> err

// Returns a list of all users
func GetAllUsers() ([]*User, error) {
	rows, err := Database.Query("SELECT * FROM `tbl_user`")
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.UserID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Domain, &parsedUser.Name)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		userArr = append(userArr, &parsedUser)
	}
	return userArr, nil
}

// Gets a user by their ID
func GetUserWithID(id uint64) (*User, error) {
	rows, err := Database.Query("SELECT * FROM `tbl_user` WHERE `id_user` = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	parsedUser := User{}
	rows.Next()
	err = rows.Scan(&parsedUser.UserID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Domain, &parsedUser.Name)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	return &parsedUser, nil
}

// Gets a user by their email
func GetUserWithEmail(email string) (*User, error) {
	rows, err := Database.Query("SELECT * FROM `tbl_user` WHERE `email` = ?", email)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	parsedUser := User{}
	rows.Next()
	err = rows.Scan(&parsedUser.UserID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Domain, &parsedUser.Name)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	return &parsedUser, nil
}

// Adds a new user to the database
func InsertUser(user *User) error {
	// Check email isn't already in use
	_, err := GetUserWithEmail(user.Email)
	if err == nil {
		fmt.Println(" [ERROR] Email '" + user.Email + "' already registered.")
		return EmailInUseError(user.Email)
	}

	_, err = Database.Query("INSERT INTO `tbl_user` (email, password, name) VALUES(?, SHA1(UNHEX(SHA1(?))), ?)", user.Email, user.Password, user.Name)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}

	loadClients()

	return nil
}

// Changes the details of a user
func UpdateUser(user *User) error {
	// Add arguments to the query if they aren't empty
	queryStr := "UPDATE `tbl_user` SET "
	argArr := make([]interface{}, 0)

	noArgs := true
	if user.Email != "" {
		noArgs = false
		queryStr += "`email` = ? ,"
		argArr = append(argArr, user.Email)

		// Check email isn't already in use
		_, err := GetUserWithEmail(user.Email)
		if err == nil {
			fmt.Println(" [ERROR] Email '" + user.Email + "' already registered.")
			return EmailInUseError(user.Email)
		}
	}

	if user.Password != "" {
		noArgs = false
		queryStr += "`password` = SHA1(UNHEX(SHA1(?))) ,"
		argArr = append(argArr, user.Password)
	}

	if user.Domain != "" {
		noArgs = false
		queryStr += "`domain` = ? ,"
		argArr = append(argArr, user.Domain)
	}

	if user.Name != "" {
		noArgs = false
		queryStr += "`name` = ? ,"
		argArr = append(argArr, user.Name)
	}

	if noArgs {
		return nil
	}

	queryStr = queryStr[:len(queryStr)-1]
	queryStr += "WHERE `id_user` = ?"
	argArr = append(argArr, user.UserID)

	_, err := Database.Query(queryStr, argArr...)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}

	// Reload clients, because names and passwords may have changed
	loadClients()

	return nil
}

// Permanently delete a user from the database
func DeleteUser(id uint64) error {
	_, err := Database.Query("DELETE FROM `tbl_user` WHERE `id_user` = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}
	return nil
}

/// ITEM FUNCTIONS ///
//	GetAllItems(userID) -> []Item & error
//  GetItemWithID(itemID) -> Item & error
//	InsertItem(Item) -> error
//  UpdateItem(Item) -> error
//  DeleteItem(itemID) -> error

// Gets all the items in the list for a user
func GetAllItems(userID uint64) ([]*Item, error) {

	// Get All Items
	rows, err := Database.Query("SELECT "+
		"i.id_item, i.name, i.desc, i.price, "+
		"s.id_status, s.name, s.desc, "+
		"u.id_user, u.email, u.domain, u.name "+
		"FROM `tbl_item` i "+
		"JOIN `tbl_status` s ON i.status_id = s.id_status "+
		"JOIN `tbl_user` u ON i.user_id = u.id_user "+
		"WHERE u.id_user = ?", userID)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	// Parse Items
	itemArr := make([]*Item, 0)
	for rows.Next() {
		parsedItem := Item{}
		err = rows.Scan(
			&parsedItem.ItemID,
			&parsedItem.Name,
			&parsedItem.Description,
			&parsedItem.Price,
			&parsedItem.Status.StatusID,
			&parsedItem.Status.Name,
			&parsedItem.Status.Description,
			&parsedItem.User.UserID,
			&parsedItem.User.Email,
			&parsedItem.User.Domain,
			&parsedItem.User.Name,
		)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}

		// Get Links
		linkRows, err := Database.Query("SELECT id_link, text, hyperlink FROM `tbl_link` WHERE item_id = ?", parsedItem.ItemID)
		if err != nil {
			fmt.Println(" [ERROR] Query Failed:", err)
			return nil, err
		}
		defer linkRows.Close()

		// Parse Links
		for linkRows.Next() {
			parsedLink := Link{}
			err = linkRows.Scan(
				&parsedLink.LinkID,
				&parsedLink.Text,
				&parsedLink.Hyperlink,
			)
			if err != nil {
				fmt.Println(" [ERROR] Parsing Failed:", err)
				return nil, err
			}
			parsedItem.Links = append(parsedItem.Links, parsedLink)
		}

		itemArr = append(itemArr, &parsedItem)
	}

	return itemArr, nil
}

// Gets a single item from the database by its ID
func GetItemWithID(id uint64) (*Item, error) {
	// Get All Items
	rows, err := Database.Query("SELECT "+
		"i.id_item, i.name, i.desc, i.price, "+
		"s.id_status, s.name, s.desc, "+
		"u.id_user, u.email, u.domain, u.name "+
		"FROM `tbl_item` i "+
		"JOIN `tbl_status` s ON i.status_id = s.id_status "+
		"JOIN `tbl_user` u ON i.user_id = u.id_user "+
		"WHERE i.id_item = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	// Parse Items
	parsedItem := Item{}
	rows.Next()
	err = rows.Scan(
		&parsedItem.ItemID,
		&parsedItem.Name,
		&parsedItem.Description,
		&parsedItem.Price,
		&parsedItem.Status.StatusID,
		&parsedItem.Status.Name,
		&parsedItem.Status.Description,
		&parsedItem.User.UserID,
		&parsedItem.User.Email,
		&parsedItem.User.Domain,
		&parsedItem.User.Name,
	)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	// Get Links
	linkRows, err := Database.Query("SELECT id_link, text, hyperlink FROM `tbl_link` WHERE item_id = ?", parsedItem.ItemID)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer linkRows.Close()

	// Parse Links
	for linkRows.Next() {
		parsedLink := Link{}
		err = linkRows.Scan(
			&parsedLink.LinkID,
			&parsedLink.Text,
			&parsedLink.Hyperlink,
		)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		parsedItem.Links = append(parsedItem.Links, parsedLink)
	}

	return &parsedItem, nil
}

// Inserts an item into the database
func InsertItem(item *Item) error {

	// Insert Item
	res, err := Database.Exec("INSERT INTO `tbl_item` (`name`, `desc`, price, status_id, user_id) VALUES (?, ?, ?, ?, ?)", item.Name, item.Description, item.Price, item.Status.StatusID, item.User.UserID)
	if err != nil {
		fmt.Println(" [ERROR] Query failed.", err)
		return err
	}

	// Get previous insert ID
	id64, err := res.LastInsertId()
	if err != nil {
		fmt.Println(" [ERROR] Failed to get previous insert ID.", err)
		return err
	}
	id := uint64(id64)

	// Insert Links
	for _, link := range item.Links {
		_, err = Database.Query("INSERT INTO `tbl_link` (text, hyperlink, item_id) VALUES (?, ?, ?)", link.Text, link.Hyperlink, id) // Use Last insert ID as itemID
		if err != nil {
			fmt.Println(" [ERROR] Query failed.", err)
			return err
		}
	}

	return nil
}

// Changes the details of an item
func UpdateItem(item *Item) error {
	// Add arguments to the query if they aren't empty
	queryStr := "UPDATE `tbl_item` SET "
	argArr := make([]interface{}, 0)

	noArgs := true
	if item.Name != "" {
		noArgs = false
		queryStr += "`name` = ? ,"
		argArr = append(argArr, item.Name)
	}

	if item.Description != "" {
		noArgs = false
		queryStr += "`desc` = ? ,"
		argArr = append(argArr, item.Description)
	}

	if item.Status != (Status{}) {
		noArgs = false
		queryStr += "`status_id` = ? ,"
		argArr = append(argArr, item.Status.StatusID)
	}

	if item.Price != -1 {
		noArgs = false
		queryStr += "`price` = ? ,"
		argArr = append(argArr, item.Price)
	}

	queryStr = queryStr[:len(queryStr)-1]
	queryStr += "WHERE `id_item` = ?"
	argArr = append(argArr, item.ItemID)

	if !noArgs {
		// Update Item
		_, err := Database.Query(queryStr, argArr...)
		if err != nil {
			fmt.Println(" [ERROR] Query failed:", err)
			return err
		}
	}

	// Drop all links for this item and add the ones provided
	if item.Links != nil {
		_, err := Database.Query("DELETE FROM tbl_link WHERE item_id = ?", item.ItemID)
		if err != nil {
			fmt.Println(" [ERROR] Query failed:", err)
			return err
		}

		for _, link := range item.Links {
			_, err = Database.Query("INSERT INTO `tbl_link` (text, hyperlink, item_id) VALUES (?, ?, ?)", link.Text, link.Hyperlink, item.ItemID)
			if err != nil {
				fmt.Println(" [ERROR] Query failed.", err)
				return err
			}
		}
	}

	return nil
}

// Permanently delete an item from the database
func DeleteItem(id uint64) error {
	_, err := Database.Query("DELETE FROM `tbl_item` WHERE `id_item` = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}
	return nil
}

/// MISC FUNCTIONS
// GetUsersByNameOrEmail(name) []User, error

// Returns a list of users with the provided substring in their email or name
func GetUsersByNameOrEmail(name string) ([]*User, error) {
	wildcardName := "%" + name + "%"
	rows, err := Database.Query("SELECT * FROM `tbl_user` WHERE LOWER(email) LIKE '?' OR LOWER(name) LIKE '?'", wildcardName, wildcardName)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.UserID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Domain, &parsedUser.Name)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		userArr = append(userArr, &parsedUser)
	}
	return userArr, nil
}
