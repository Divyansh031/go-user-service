package unit

import (
	"testing"
	"time"

	"github.com/Divyansh031/user-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	assert.NotNil(t, user)
	assert.NotEqual(t, "", user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, dob, user.DateOfBirth)
	assert.Equal(t, "+1234567890", user.PhoneNumber)
	assert.Equal(t, "john@example.com", user.Email)
	assert.False(t, user.IsBlocked)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name      string
		user      *domain.User
		wantError error
	}{
		{
			name: "valid user",
			user: &domain.User{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "male",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "+1234567890",
				Email:       "john@example.com",
			},
			wantError: nil,
		},
		{
			name: "missing first name",
			user: &domain.User{
				LastName:    "Doe",
				Gender:      "male",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "+1234567890",
				Email:       "john@example.com",
			},
			wantError: domain.ErrInvalidFirstName,
		},
		{
			name: "missing last name",
			user: &domain.User{
				FirstName:   "John",
				Gender:      "male",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "+1234567890",
				Email:       "john@example.com",
			},
			wantError: domain.ErrInvalidLastName,
		},
		{
			name: "invalid gender",
			user: &domain.User{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "invalid",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "+1234567890",
				Email:       "john@example.com",
			},
			wantError: domain.ErrInvalidGender,
		},
		{
			name: "future date of birth",
			user: &domain.User{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "male",
				DateOfBirth: time.Now().Add(24 * time.Hour),
				PhoneNumber: "+1234567890",
				Email:       "john@example.com",
			},
			wantError: domain.ErrInvalidDateOfBirth,
		},
		{
			name: "missing phone number",
			user: &domain.User{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "male",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				Email:       "john@example.com",
			},
			wantError: domain.ErrInvalidPhoneNumber,
		},
		{
			name: "missing email",
			user: &domain.User{
				FirstName:   "John",
				LastName:    "Doe",
				Gender:      "male",
				DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "+1234567890",
			},
			wantError: domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(time.Millisecond * 10) // Ensure time difference

	newDob := time.Date(1991, 2, 20, 0, 0, 0, 0, time.UTC)
	user.Update("Jane", "Smith", "female", newDob)

	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, "female", user.Gender)
	assert.Equal(t, newDob, user.DateOfBirth)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestUserBlock(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	assert.False(t, user.IsBlocked)

	user.Block()
	assert.True(t, user.IsBlocked)
}

func TestUserUnblock(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	user.Block()
	assert.True(t, user.IsBlocked)

	user.Unblock()
	assert.False(t, user.IsBlocked)
}

func TestUserUpdateContact(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	newPhone := "+9876543210"
	newEmail := "newemail@example.com"

	user.UpdateContact(&newPhone, &newEmail)

	assert.Equal(t, newPhone, user.PhoneNumber)
	assert.Equal(t, newEmail, user.Email)
}

func TestUserUpdateContactPartial(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := domain.NewUser("John", "Doe", "male", dob, "+1234567890", "john@example.com")

	originalPhone := user.PhoneNumber
	newEmail := "newemail@example.com"

	// Update only email
	user.UpdateContact(nil, &newEmail)

	assert.Equal(t, originalPhone, user.PhoneNumber) // Phone unchanged
	assert.Equal(t, newEmail, user.Email)            // Email updated
}