package validator

import (
	"regexp"
)

var (
	// Email regex pattern
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Phone regex pattern (E.164 format: +[1-9]{1}[0-9]{1,14})
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	if len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidatePhoneNumber validates phone number format (E.164)
func ValidatePhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

// ValidateFirstName validates first name
func ValidateFirstName(firstName string) bool {
	if firstName == "" {
		return false
	}
	if len(firstName) < 2 || len(firstName) > 50 {
		return false
	}
	return true
}

// ValidateLastName validates last name
func ValidateLastName(lastName string) bool {
	if lastName == "" {
		return false
	}
	if len(lastName) < 2 || len(lastName) > 50 {
		return false
	}
	return true
}

// ValidateGender validates gender
func ValidateGender(gender string) bool {
	validGenders := map[string]bool{
		"male":   true,
		"female": true,
		"other":  true,
	}
	return validGenders[gender]
}
