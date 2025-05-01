package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/core"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username       string    `gorm:"uniqueIndex" json:"username"`
	Email          string    `gorm:"uniqueIndex" json:"-"`
	HashedPassword string    `gorm:"column:hashed_password" json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewUser(username ...string) *User {
	if len(username) == 0 {
		return &User{}
	}

	return &User{
		Username: username[0],
	}
}

func (u *User) WithID(id uuid.UUID) *User {
	u.ID = id
	return u
}

func (u *User) WithUsername(username string) *User {
	u.Username = username
	return u
}

func (u *User) WithEmail(email string) *User {
	u.Email = email
	return u
}

// Automatically hashes the password
func (u *User) WithPassword(unhashedPassword string) *User {
	hashedPassword, salt, _ := core.CreateHashWithSalt(unhashedPassword)
	u.HashedPassword = core.SerializeHashWithSalt(hashedPassword, salt)
	return u
}

func (u *User) Create(db *gorm.DB) error {
	result := db.Create(u)
	return result.Error
}

// FindByUsername looks up a user by username
func (u *User) FindByUsername(db *gorm.DB, username ...string) error {
	var _username string

	if len(username) > 0 {
		_username = username[0]
	} else {
		_username = u.Username
	}

	result := db.Where("username = ?", _username).First(u)
	return result.Error
}

// FindByEmail looks up a user by email
func (u *User) FindByEmail(db *gorm.DB, email ...string) error {
	var _email string

	if len(email) > 0 {
		_email = email[0]
	} else {
		_email = u.Email
	}

	result := db.Where("email = ?", _email).First(u)
	return result.Error
}

func (u *User) FindByID(db *gorm.DB, id ...uuid.UUID) error {
	var _id uuid.UUID

	if len(id) > 0 {
		_id = id[0]
	} else {
		_id = u.ID
	}

	result := db.First(u, _id)
	return result.Error
}

func (u *User) Update(db *gorm.DB) error {
	result := db.Save(u)
	return result.Error
}

func (u *User) Delete(db *gorm.DB) error {
	result := db.Delete(u)
	return result.Error
}
