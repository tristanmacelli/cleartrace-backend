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

//GetByID gogog
func (ms *MysqlStore) GetByID(id int64) (*User, error) {
	// open a connection

	queryString := "SELECT * FROM users where ID = " + strconv.FormatInt(id, 10)
	rows, err := ms.db.Query(queryString)
	if err != nil {
		// close the db connefction
		return nil, err
	}
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
			// close the db connefction
			return nil, err
		}
	}
	// close the db connefction
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

// type User struct {
// 	ID        int64  `json:"id"`
// 	Email     string `json:"-"` //never JSON encoded/decoded
// 	PassHash  []byte `json:"-"` //never JSON encoded/decoded
// 	UserName  string `json:"userName"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	PhotoURL  string `json:"photoURL"`
// }
