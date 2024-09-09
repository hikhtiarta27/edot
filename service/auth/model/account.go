package model

import (
	"net/mail"
	"shared"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	ID        ulid.ULID
	Name      string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (Account) TableName() string {
	return "accounts"
}

func NewAccount(
	name string,
	username string,
	password string,
) (*Account, error) {
	instance := &Account{
		ID:   ulid.Make(),
		Name: strings.TrimSpace(name),
	}

	if err := instance.SetUsername(username); err != nil {
		return nil, err
	}

	if err := instance.SetPassword(password); err != nil {
		return nil, err
	}

	return instance, nil
}

func (m *Account) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}

	m.Password = string(hashedPassword)
	return nil
}

func (m *Account) SetUsername(username string) error {
	username = strings.TrimSpace(username)

	if _, err := mail.ParseAddress(username); err != nil {
		username = m.parsePhoneNumber(username)
	}

	m.Username = username
	return nil
}

func (m Account) parsePhoneNumber(phoneNumber string) string {

	phoneNumber = strings.TrimSpace(phoneNumber)

	if phoneNumber[0] == '+' {
		phoneNumber = phoneNumber[1:]
	}

	return phoneNumber
}

func (m Account) ComparePassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(plainPassword))
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 401,
			Message:        "wrong username/password",
		}

	}

	return nil
}
