package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the domain
type User struct {
	ID          string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Email       string
	IsBlocked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewUser creates a new user with generated ID and timestamps
func NewUser(firstName, lastName, gender string, dob time.Time, phone, email string) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New().String(),
		FirstName:   firstName,
		LastName:    lastName,
		Gender:      gender,
		DateOfBirth: dob,
		PhoneNumber: phone,
		Email:       email,
		IsBlocked:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates user fields
func (u *User) Update(firstName, lastName, gender string, dob time.Time) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Gender = gender
	u.DateOfBirth = dob
	u.UpdatedAt = time.Now()
}

// UpdateContact updates contact information
func (u *User) UpdateContact(phone, email *string) {
	if phone != nil {
		u.PhoneNumber = *phone
	}
	if email != nil {
		u.Email = *email
	}
	u.UpdatedAt = time.Now()
}

// Block blocks the user
func (u *User) Block() {
	u.IsBlocked = true
	u.UpdatedAt = time.Now()
}

// Unblock unblocks the user
func (u *User) Unblock() {
	u.IsBlocked = false
	u.UpdatedAt = time.Now()
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.FirstName == "" {
		return ErrInvalidFirstName
	}
	if u.LastName == "" {
		return ErrInvalidLastName
	}
	if u.Gender == "" || !isValidGender(u.Gender) {
		return ErrInvalidGender
	}
	if u.DateOfBirth.IsZero() || u.DateOfBirth.After(time.Now()) {
		return ErrInvalidDateOfBirth
	}
	if u.PhoneNumber == "" {
		return ErrInvalidPhoneNumber
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

func isValidGender(gender string) bool {
	validGenders := []string{"male", "female", "other"}
	for _, g := range validGenders {
		if gender == g {
			return true
		}
	}
	return false
}