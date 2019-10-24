package users

import (
	"fmt"
	"strconv"
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

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByID(1); err != nil {
		t.Errorf("was not expecting an error, but there was none")
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

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("user@domain.com").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByEmail("user@domain.com"); err != nil {
		t.Errorf("was not expecting an error, but there was none")
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

	mock.ExpectPrepare("SELECT \\* FROM users")
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("Sam").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.GetByUserName("Sam"); err != nil {
		t.Errorf("was not expecting an error, but there was none")
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

	nu := NewUser{"user@domain.com", "password", "password", "username", "first", "last"}
	u, err := nu.ToUser()
	if err != nil {
		fmt.Println("Error when generating user")
	}
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users\\(email, passHash, username, firstname, lastname, photoURL\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if user, err := ms.Insert(u); err != nil {
		t.Errorf("was not expecting an error, but there was none")
		fmt.Print(user)
	}
}

func TestInsertExpectError(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	nu := NewUser{"user@domain.com", "password", "password", "username", "first", "last"}
	u, err := nu.ToUser()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users\\(email, passHash, username firstname, lastname, photoURL\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnError(fmt.Errorf("Some error"))
	mock.ExpectRollback()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if _, err := ms.Insert(u); err == nil {
		t.Errorf("we were expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = 1
	updates := &Updates{"testfirst", "testlast"}
	columns := []string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE users SET firstname = \\?, lastname = \\? WHERE ID = \\?").
		WithArgs(updates.FirstName, updates.LastName, strconv.FormatInt(id, 10)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if _, err := ms.Update(id, updates); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateExpectError(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = -1
	updates := &Updates{"testfirst", "testlast"}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE users SET firstname = \\?, lastname = \\? WHERE ID = \\?").
		WithArgs(updates.FirstName, updates.LastName, strconv.FormatInt(id, 10)).
		WillReturnError(fmt.Errorf("No negative ID values"))
	mock.ExpectRollback()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if _, err := ms.Update(id, updates); err == nil {
		t.Errorf("we were expecting an error, but there were none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = 1

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users").
		WithArgs(strconv.FormatInt(id, 10)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if err := ms.Delete(id); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

//TODO: Write TestDeleteExpectError
func TestDeleteExpectError(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = -1

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users").
		WithArgs(strconv.FormatInt(id, 10)).
		WillReturnError(fmt.Errorf("No negative ID values"))
	mock.ExpectRollback()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.db = db

	// now we execute our method with the mock
	if err := ms.Delete(id); err == nil {
		t.Errorf("we were expecting an error, but there were none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
