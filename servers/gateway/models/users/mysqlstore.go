package users

import (
	"database/sql"
	"fmt"
	"os"
)

//RedisStore represents a session.Store backed by redis.
type MysqlStore struct {
}

func NewMysqlStore() *MysqlStore {
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/demo", os.Getenv("MYSQL_ROOT_PASSWORD"))

	//create a database object, which manages a pool of
	//network connections to the database server
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return nil
	}

}
