package lib

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	password []byte
}

func (a *Admin) ValidatePassword(userInput string) bool {
	err := bcrypt.CompareHashAndPassword(a.password, []byte(userInput))
	if err != nil {
		fmt.Println("Invalid Password: ", err)
		return false
	}
	fmt.Println("Login successful.")
	return true
}

func (a *Admin) Exists() bool {
	var value []byte
	err := db.QueryRow("SELECT value FROM settings WHERE name = 'password';").Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No Admin-Password found.")
		return false
	case err != nil:
		fmt.Println("Error while checking for Admin-Password: ", err)
		return false
	default:
		a.password = value
	}

	fmt.Println("Existing Admin-Password found.")
	return true
}

func (a *Admin) ChangePassword(userInput string) {
	fmt.Println("Changing Admin-Password...")

	clearPassword := []byte(userInput)
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 15)

	stmt, err := db.Prepare("UPDATE settings SET value = (?) WHERE name = 'password';")
	if err != nil {
		fmt.Println("Unable to update Admin-Password: ", err)
	} else {
		_, err = stmt.Exec(passwordHash)
		if err != nil {
			fmt.Println("Unable to update Admin-Password: ", err)
		}
	}

	a.password = passwordHash
}

func (a *Admin) Add() {
	fmt.Println("Adding default Admin...")

	clearPassword := []byte("admin")
	passwordHash, err := bcrypt.GenerateFromPassword(clearPassword, 15)

	stmt, err := db.Prepare("INSERT INTO settings (name, value) VALUES ('password', (?));")
	if err != nil {
		fmt.Println("Unable to insert Admin-Password: ", err)
	} else {
		_, err = stmt.Exec(passwordHash)
		if err != nil {
			fmt.Println("Unable to insert Admin-Password: ", err)
		}
	}

	a.password = passwordHash
}
