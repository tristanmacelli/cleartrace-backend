package users

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// a failing test case
func TestGetByID(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByID(1); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByEmail(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("user@domain.com").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByEmail("user@domain.com"); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByUsername(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("Sam").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByUserName("Sam"); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	nu := NewUser{}
	nu.Email = "user@domain.com"
	nu.Password = "password"
	nu.PasswordConf = "password"
	nu.UserName = "username"
	nu.FirstName = "first"
	nu.LastName = "last"
	u, err := nu.ToUser()
	if err != nil {
		fmt.Println("Error when generating user")
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users\\(email, passHash, username, firstname, lastname, photoURL\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if _, err := ms.Insert(u); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
