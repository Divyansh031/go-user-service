package unit

import (
	"testing"

	"github.com/Divyansh031/user-service/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "valid email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "valid email with numbers",
			email:    "user123@example.com",
			expected: true,
		},
		{
			name:     "valid email with special chars",
			email:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "empty email",
			email:    "",
			expected: false,
		},
		{
			name:     "email without @",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "email without domain",
			email:    "user@",
			expected: false,
		},
		{
			name:     "email with spaces",
			email:    "user @example.com",
			expected: false,
		},
		{
			name:     "email too long",
			email:    "a" + string(make([]byte, 255)) + "@example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{
			name:     "valid phone US format",
			phone:    "+12015550123",
			expected: true,
		},
		{
			name:     "valid phone IN format",
			phone:    "+919876543210",
			expected: true,
		},
		{
			name:     "valid phone international",
			phone:    "+447911123456",
			expected: true,
		},
		{
			name:     "empty phone",
			phone:    "",
			expected: false,
		},
		{
			name:     "phone without plus",
			phone:    "12015550123",
			expected: false,
		},
		{
			name:     "phone with spaces",
			phone:    "+1 201 555 0123",
			expected: false,
		},
		{
			name:     "phone too short",
			phone:    "+1",
			expected: false,
		},
		{
			name:     "phone starting with +0",
			phone:    "+01234567890",
			expected: false,
		},
		{
			name:     "phone with letters",
			phone:    "+1201555ABCD",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePhoneNumber(tt.phone)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateFirstName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		expected  bool
	}{
		{
			name:      "valid name",
			firstName: "John",
			expected:  true,
		},
		{
			name:      "valid long name",
			firstName: "Alexanderthe",
			expected:  true,
		},
		{
			name:      "empty name",
			firstName: "",
			expected:  false,
		},
		{
			name:      "too short",
			firstName: "J",
			expected:  false,
		},
		{
			name:      "too long",
			firstName: string(make([]byte, 51)),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateFirstName(tt.firstName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateLastName(t *testing.T) {
	tests := []struct {
		name     string
		lastName string
		expected bool
	}{
		{
			name:     "valid name",
			lastName: "Doe",
			expected: true,
		},
		{
			name:     "valid long name",
			lastName: "VanDerberg",
			expected: true,
		},
		{
			name:     "empty name",
			lastName: "",
			expected: false,
		},
		{
			name:     "too short",
			lastName: "D",
			expected: false,
		},
		{
			name:     "too long",
			lastName: string(make([]byte, 51)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateLastName(tt.lastName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateGender(t *testing.T) {
	tests := []struct {
		name     string
		gender   string
		expected bool
	}{
		{
			name:     "valid male",
			gender:   "male",
			expected: true,
		},
		{
			name:     "valid female",
			gender:   "female",
			expected: true,
		},
		{
			name:     "valid other",
			gender:   "other",
			expected: true,
		},
		{
			name:     "empty gender",
			gender:   "",
			expected: false,
		},
		{
			name:     "invalid gender",
			gender:   "unknown",
			expected: false,
		},
		{
			name:     "case sensitive",
			gender:   "MALE",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateGender(tt.gender)
			assert.Equal(t, tt.expected, result)
		})
	}
}
