package domain

import "errors"

// Error variables for domain errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
	ErrInvalidFirstName   = errors.New("invalid first name")
	ErrInvalidLastName    = errors.New("invalid last name")
	ErrInvalidGender      = errors.New("invalid gender")
	ErrInvalidDateOfBirth = errors.New("invalid date of birth")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrDatabaseError      = errors.New("database error")
	ErrInternal           = errors.New("internal server error")
)
