package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	// Necessary to run db commands
	_ "github.com/go-sql-driver/mysql"
)

//MysqlStore represents a connection to our user database
type MysqlStore struct {
	DB *sql.DB
}

// A partially constructed sql query to use in getter functions
const queryString = "SELECT * FROM users WHERE"

//NewMysqlStore creates an open database connection to do queries and transactions on
// For help on forming dsn (https://drstearns.github.io/tutorials/godb/#secconnectingfromagoprogram)
// See docker run command for env vars that define database name & password
func NewMysqlStore(dsn string) *MysqlStore {
	// We are using a persistent connection for all transactions
	db, err := sql.Open("mysql", dsn)
	fmt.Println("Error opening db", err)
	return &MysqlStore{
		DB: db,
	}
}

// GetBy is a helper method that resolves all queries asking for singular user objects
func (ms *MysqlStore) GetBy(query string, value string) (*User, error) {
	// Creates a new user object to populate with the query for a single row
	user := User{}
	insq := queryString + query
	row := ms.DB.QueryRow(insq, value)
	// Populating the new user
	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName, &user.PhotoURL)

	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			// We return a user without any values indicating that the query returned nothing
			// since there was no fatal error
			fmt.Printf("error getting user: %v\n", err)
		} else {
			return nil, err
		}
	}
	return &user, nil
}

//GetByID returns the User with the given ID
func (ms *MysqlStore) GetByID(id int64) (*User, error) {
	query := " ID = ?"
	return ms.GetBy(query, strconv.FormatInt(id, 10))
}

//GetByEmail returns the User with the given email
func (ms *MysqlStore) GetByEmail(email string) (*User, error) {
	query := " Email = ?"
	return ms.GetBy(query, email)
}

//GetByUserName returns the User with the given Username
func (ms *MysqlStore) GetByUserName(username string) (*User, error) {
	query := " UserName = ?"
	return ms.GetBy(query, username)
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MysqlStore) Insert(user *User) (*User, error) {
	// This inserts a new row into the "users" table Using ? markers for the values will defeat SQL
	// injection attacks
	fmt.Println(user.Email)
	fmt.Println(user.PassHash)
	fmt.Println(user.UserName)
	fmt.Println(user.FirstName)
	fmt.Println(user.LastName)
	fmt.Println(user.PhotoURL)
	insq := "INSERT INTO users(email, passHash, username, firstname, lastname, photoURL) VALUES (?,?,?,?,?,?)"
	res, err := ms.DB.Exec(insq, user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return nil, err
	}
	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error getting new ID: %v\n", id)
		return nil, err
	}
	fmt.Printf("ID for new row is %d\n", id)
	// Get and return this new user
	return ms.GetByID(id)
}

//LogSuccessfulSignIns does something
func (ms *MysqlStore) LogSuccessfulSignIns(user *User, r *http.Request) {
	uid := user.ID
	timeOfSignIn := time.Now()
	clientIP := r.RemoteAddr
	ips := r.Header.Get("X-Forwarded-For")

	if len(ips) > 1 {
		clientIP = strings.Split(ips, ",")[0]
	} else if len(ips) == 1 {
		clientIP = ips
	}
	insq := "INSERT INTO userSignIn(userID, signinDT, ip) VALUES (?,?,?)"
	_, err := ms.DB.Exec(insq, uid, timeOfSignIn, clientIP)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return
	}
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MysqlStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "UPDATE users SET firstname = ?, lastname = ? WHERE ID = ?"
	_, err := ms.DB.Exec(insq, updates.FirstName, updates.LastName, strconv.FormatInt(id, 10))
	if err != nil {
		fmt.Printf("error updating row: %v\n", err)
		return nil, err
	}
	return ms.GetByID(id)
}

//Delete deletes the user with the given ID
func (ms *MysqlStore) Delete(id int64) error {
	insq := "DELETE FROM users WHERE ID = ?"
	_, err := ms.DB.Exec(insq, strconv.FormatInt(id, 10))
	if err != nil {
		fmt.Printf("error deleting row: %v\n", err)
		return err
	}
	return nil
}
