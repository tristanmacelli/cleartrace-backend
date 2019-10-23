package users

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//RedisStore represents a session.Store backed by redis.
type MysqlStore struct {
	dsn string
	db  DB
}

func NewMysqlStore() *MysqlStore {
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/demo", os.Getenv("MYSQL_ROOT_PASSWORD"))
	return &MysqlStore{
		dsn: dsn,
	}
}

func OpenConnection(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return nil, err
	}
	return db, nil
}

func (ms *MysqlStore) GetByID(id int64) (*User, error) {
	// open a connection

	queryString := "SELECT * FROM users where ID = " + id
	rows, err := db.Query(queryString)
	if err != nil {
		// close the fdb connefction
		return nil, err
	}
	var ID int
	var Email string
	var PassHash []byte
	var UserName string
	var FirstName string
	var LastName string
	var PhotoURL string
	for rows.Next() {
		err = rows.Scan(&ID, &Email, &PassHash, &UserName, &FirstName, &LastName, &PhotoURL)
		if err != nil {
			// close the fdb connefction
			return nil, err
		}
	}
	// close the fdb connefction
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
