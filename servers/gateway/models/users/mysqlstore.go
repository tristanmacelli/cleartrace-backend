package users

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	// Necessary to run db commands
	_ "github.com/go-sql-driver/mysql"
)

//MysqlStore represents a connection to our user database
type MysqlStore struct {
	db *sql.DB
}

// A partially constructed sql query to use in getter functions
const queryString = "SELECT * FROM users WHERE"

//NewMysqlStore creates data source name which can be used to connect to the user database
func NewMysqlStore() *MysqlStore {
	// See docker run command for env vars that define database name & password
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/insert-database-name-here", os.Getenv("MYSQL_ROOT_PASSWORD"))
	// We are using a persistent connection for all transactions
	db, _ := sql.Open("mysql", dsn)
	return &MysqlStore{
		db: db,
	}
}

// GetBy is a helper method that resolves all queries asking for singular user objects
func (ms *MysqlStore) GetBy(query string, queryValue string) (*User, error) {

	// A prepared statement allows us to use bind values the (?)
	// in the query parameter which helps avoid sql injection attacks
	stmt, _ := ms.db.Prepare(queryString + query)
	defer stmt.Close()
	row := stmt.QueryRow(queryValue)
	// May need to alter the capitalization of these variables to match what is in
	// schema.sql or vice versa
	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
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
	// This inserts a new row into the "users" table
	// Using ? markers for the values will defeat SQL
	// injection attacks

	// Open a reserved connection to make an individual transaction
	tx, _ := ms.db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO users(email, passHash, username, firstname, lastname, photoURL) VALUES (?,?,?,?,?,?)")
	defer stmt.Close()
	res, err := stmt.Exec(user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		// Close the reserved connection upon failure
		tx.Rollback()
		return nil, err
	}
	// Close the reserved connection upon success
	tx.Commit()

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

// //Update applies UserUpdates to the given user ID
// //and returns the newly-updated user
// func (ms *MysqlStore) Update(id int64, updates *Updates) (*User, error) {

// 	// Open a reserved connection to db to make an individual transaction
// 	tx, _ := ms.db.Begin()
// 	stmt, _ := tx.Prepare("UPDATE users SET firstname = ?, lastname = ? WHERE ID = ?")
// 	// This will close the prepared statement once Exec is called
// 	defer stmt.Close()
// 	res, err := stmt.Exec(updates.FirstName, updates.LastName, strconv.FormatInt(id, 10))
// 	if err != nil {
// 		fmt.Printf("error updating row: %v\n", err)
// 		// Close the reserved connection upon failure
// 		tx.Rollback()
// 		return nil, err
// 	}
// 	// Close the reserved connection upon success
// 	tx.Commit()
// 	return nil, nil
// }

//Delete deletes the user with the given ID
func (ms *MysqlStore) Delete(id int64) error {

	// Open a reserved connection to db to make an individual transaction
	tx, _ := ms.db.Begin()
	stmt, _ := tx.Prepare("DELETE FROM users WHERE ID = ?")
	// This will close the prepared statement once Exec is called
	defer stmt.Close()
	_, err := stmt.Exec(strconv.FormatInt(id, 10))
	if err != nil {
		fmt.Printf("error deleting row: %v\n", err)
		// Close the reserved connection upon failure
		tx.Rollback()
		return err
	}
	// Close the reserved connection upon success
	tx.Commit()
	return nil
}
