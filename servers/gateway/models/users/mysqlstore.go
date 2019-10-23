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
	dsn string
	db  *sql.DB
}

// A partially constructed sql query to use in getter functions
const queryString = "SELECT * FROM users where"

//NewMysqlStore creates data source name which can be used to connect to the user database
func NewMysqlStore() *MysqlStore {
	// See docker run command for env vars that define database name & password
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/insert-database-name-here", os.Getenv("MYSQL_ROOT_PASSWORD"))
	// Uncomment this if we can use a persistent connection for all transactions
	// db, _ := sql.Open("mysql", dsn)
	return &MysqlStore{
		dsn: dsn,
		// Uncomment this if we can use a persistent connection for all transactions
		// db:  db,
	}
}

//OpenConnection opens a connection to the user database
// Delete this function if we can use a persistent connection for all transactions
//	i.e. we are storing the open connection in our MysqlStore struct
// Keep this function if we must use a different connection for each individual transactions
//  i.e. we are opening a connection in each of our CRUD functions
func OpenConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return nil, err
	}
	return db, nil
}

// GetBy is a helper method that resolves all
func (ms *MysqlStore) GetBy(query string) (*User, error) {

	// TODO: Open connection to db
	rows, err := ms.db.Query(query)
	// TODO: Close connection to db
	if err != nil {
		return nil, err
	}
	// May need to alter the capitalization of these variables to match what is in
	// schema.sql or vice versa
	var ID int64
	var Email string
	var PassHash []byte
	var UserName string
	var FirstName string
	var LastName string
	var PhotoURL string
	for rows.Next() {
		err = rows.Scan(&ID, &Email, &PassHash, &UserName, &FirstName, &LastName, &PhotoURL)
		if err != nil {
			return nil, err
		}
	}
	return &User{
		ID:        ID,
		Email:     Email,
		PassHash:  PassHash,
		UserName:  UserName,
		FirstName: FirstName,
		LastName:  LastName,
		PhotoURL:  PhotoURL,
	}, nil
}

//GetByID returns the User with the given ID
func (ms *MysqlStore) GetByID(id int64) (*User, error) {
	query := queryString + " ID = " + strconv.FormatInt(id, 10)
	return ms.GetBy(query)
}

//GetByEmail returns the User with the given email
func (ms *MysqlStore) GetByEmail(email string) (*User, error) {
	query := queryString + " Email = " + email
	return ms.GetBy(query)
}

//GetByUserName returns the User with the given Username
func (ms *MysqlStore) GetByUserName(username string) (*User, error) {
	query := queryString + " UserName = " + username
	return ms.GetBy(query)
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MysqlStore) Insert(user *User) (*User, error) {
	//insert a new row into the "users" table
	//use ? markers for the values to defeat SQL
	//injection attacks

	// TODO: Open connection to db
	insq := "insert into users(email, passHash, username, firstname, lastname, photoURL) values (?,?,?,?,?,?)"
	res, err := ms.db.Exec(insq, user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)
	// TODO: Close connection to db

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
	return ms.GetByID(id)
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MysqlStore) Update(id int64, updates *Updates) (*User, error) {
	// TODO: Implement this function per the comment above
	return nil, nil
}

//Delete deletes the user with the given ID
func (ms *MysqlStore) Delete(id int64) error {
	// TODO: Implement this function per the comment above
	return nil
}

// type User struct {
// 	ID        int64  `json:"id"`
// 	Email     string `json:"-"` //never JSON encoded/decoded
// 	PassHash  []byte `json:"-"` //never JSON encoded/decoded
// 	UserName  string `json:"userName"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	PhotoURL  string `json:"photoURL"`
// }
