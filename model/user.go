package model

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// User ...
type User struct {
	// ID ...
	ID bson.ObjectId `json:"id" bson:"_id"`
	// Username ...
	Username string `json:"username" bson:"username"`
	// Email ...
	Email string `json:"email" bson:"email"`
	// PasswordHash ...
	PasswordHash string `json:"-" bson:"password_hash"`
}

// GeneratePasswordHash ...
func (u *User) GeneratePasswordHash(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// ValidatePassword ...
func (u User) ValidatePassword(password string) bool {
	hashedPassword := []byte(u.PasswordHash)
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
