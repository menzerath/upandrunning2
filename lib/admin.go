package lib

import (
	"database/sql"
	"github.com/op/go-logging"
	"golang.org/x/crypto/bcrypt"
)

// Contains the current password hash.
type Admin struct {
	password []byte
}

// Init an Admin-struct.
// Creates an Admin user if there is none.
func (a *Admin) Init() {
	if !a.LoadPassword() {
		a.Add()
	}
}

// Validates the entered password.
// Returns true, if the password matched the stored hash and false, if not.
func (a *Admin) ValidatePassword(userInput string) bool {
	a.LoadPassword()
	err := bcrypt.CompareHashAndPassword(a.password, []byte(userInput))
	if err != nil {
		logging.MustGetLogger("").Warning("Invalid Password: ", err)
		return false
	}
	logging.MustGetLogger("").Info("Login successful.")
	return true
}

// Loads the current password hash into the Admin-struct.
// Returns true, if there was a password hash in the database or false, if not.
func (a *Admin) LoadPassword() bool {
	var value []byte
	err := db.QueryRow("SELECT value FROM settings WHERE name = 'password';").Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		logging.MustGetLogger("").Warning("No Admin-Password found.")
		return false
	case err != nil:
		logging.MustGetLogger("").Error("Error while checking for Admin-Password: ", err)
		return false
	default:
		a.password = value
	}

	logging.MustGetLogger("").Debug("Existing Admin-Password found.")
	return true
}

// Changes the current password to the given one.
// Returns an error (if there was one).
func (a *Admin) ChangePassword(userInput string) error {
	logging.MustGetLogger("").Debug("Changing Admin-Password...")

	clearPassword := []byte(userInput)
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 12)

	_, err = db.Exec("UPDATE settings SET value = ? WHERE name = 'password';", passwordHash)
	if err != nil {
		logging.MustGetLogger("").Error("Unable to update Admin-Password: ", err)
	}

	a.password = passwordHash
	return err
}

// Adds a new admin user to the database.
func (a *Admin) Add() {
	logging.MustGetLogger("").Info("Adding default Admin...")

	clearPassword := []byte("admin")
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 12)

	_, err = db.Exec("INSERT INTO settings (name, value) VALUES ('password', ?);", passwordHash)
	if err != nil {
		logging.MustGetLogger("").Fatal("Unable to insert Admin-Password: ", err)
	}

	a.password = passwordHash
}
