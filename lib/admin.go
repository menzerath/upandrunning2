package lib

import (
	"database/sql"
	"github.com/op/go-logging"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	password []byte
}

func (a *Admin) Init() {
	if !a.LoadPassword() {
		a.Add()
	}
}

func (a *Admin) ValidatePassword(userInput string) bool {
	a.LoadPassword()
	err := bcrypt.CompareHashAndPassword(a.password, []byte(userInput))
	if err != nil {
		logging.MustGetLogger("logger").Warning("Invalid Password: ", err)
		return false
	}
	logging.MustGetLogger("logger").Info("Login successful.")
	return true
}

func (a *Admin) LoadPassword() bool {
	var value []byte
	err := db.QueryRow("SELECT value FROM settings WHERE name = 'password';").Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		logging.MustGetLogger("logger").Warning("No Admin-Password found.")
		return false
	case err != nil:
		logging.MustGetLogger("logger").Error("Error while checking for Admin-Password: ", err)
		return false
	default:
		a.password = value
	}

	logging.MustGetLogger("logger").Debug("Existing Admin-Password found.")
	return true
}

func (a *Admin) ChangePassword(userInput string) error {
	logging.MustGetLogger("logger").Debug("Changing Admin-Password...")

	clearPassword := []byte(userInput)
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 15)

	stmt, err := db.Prepare("UPDATE settings SET value = (?) WHERE name = 'password';")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to update Admin-Password: ", err)
	} else {
		_, err = stmt.Exec(passwordHash)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to update Admin-Password: ", err)
		}
	}

	a.password = passwordHash
	return err
}

func (a *Admin) Add() {
	logging.MustGetLogger("logger").Info("Adding default Admin...")

	clearPassword := []byte("admin")
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 15)

	stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('password', (?));")
	if err != nil {
		logging.MustGetLogger("logger").Error("Unable to insert Admin-Password: ", err)
	} else {
		_, err = stmt.Exec(passwordHash)
		if err != nil {
			logging.MustGetLogger("logger").Error("Unable to insert Admin-Password: ", err)
		}
	}

	a.password = passwordHash
}
