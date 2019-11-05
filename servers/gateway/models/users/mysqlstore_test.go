package users

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID tests to see that we can query the database with a valid ID
// and receive the all the column data corresponding to the row with that ID without errors
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
	ms.DB = db

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

// TestGetByIDExpectError tests to see that when we query the database with an invalid ID
// and we do not receive column data corresponding any row and that we return an error
func TestGetByIDExpectError(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("-1").
		WillReturnError(fmt.Errorf("Invalid ID value: ID's must be positive"))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.DB = db

	// now we execute our method with the mock
	if user, err := ms.GetByID(-1); err == nil {
		t.Errorf("we were expecting an error, but there was none")
		fmt.Print(user)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetByEmail tests to see that we can query the database with a valid email
// and receive the all the column data corresponding to the row with that email without errors
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
	ms.DB = db

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

// TestGetByUsername tests to see that we can query the database with a valid username
// and receive the all the column data corresponding to the row with that username without errors
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
	ms.DB = db

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

// TestInsert tests to see that upon calling the update method with a validated user,
// the user will be committed into the database without errors and then the new user will
// be made available
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
	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs("1").
		WillReturnRows(mock.NewRows(columns))

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.DB = db

	// now we execute our method with the mock
	if _, err := ms.Insert(u); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestExpectErrorInsert Tests case when a user with an undefined passHash is trying to
// be inserted. Expect an error and a transaction rollback (do not commit new data to
// the database)
func TestExpectErrorInsert(t *testing.T) {
	//MysqlStore represents a connection to our user database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	nu := NewUser{"user@domain.com", "password", "password", "username", "first", "last"}
	u, err := nu.ToUser()
	u.PassHash = []byte{}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).
		WillReturnError(fmt.Errorf("No password passed"))
	mock.ExpectRollback()

	// passes the mock to our struct
	var ms = MysqlStore{}
	ms.DB = db

	// now we execute our method with the mock
	if _, err := ms.Insert(u); err == nil {
		t.Errorf("we were expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdate tests the update operation in the Update function as well as checking that
// the data store returns the updated user after the update transaction was committed sucessfully
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
	ms.DB = db

	// now we execute our method with the mock
	if _, err := ms.Update(id, updates); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateExpectError tests an update using an invalid row id, which should cause
// the update transaction to fail and rolling back the changes.
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
	ms.DB = db

	// now we execute our method with the mock
	if _, err := ms.Update(id, updates); err == nil {
		t.Errorf("we were expecting an error, but there were none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDelete tests to see that, given a valid user ID, a user will be deleted from
// the database and this change will be committed to the database without error.
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
	ms.DB = db

	// now we execute our method with the mock
	if err := ms.Delete(id); err != nil {
		t.Errorf("there was an error, but we were not expecting one")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteExpectError attempts to delete a user with an invalid ID which causes
// an error and causes the transaction to rollback (no data deleted)
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
	ms.DB = db

	// now we execute our method with the mock
	if err := ms.Delete(id); err == nil {
		t.Errorf("we were expecting an error, but there were none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
