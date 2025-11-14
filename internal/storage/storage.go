// internal/storage/storage.go
package storage

import (
	"context"

	"github.com/Divyansh031/user-service/internal/domain"
	"github.com/google/uuid"
)

type Storage interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit int, pageToken string) ([]*domain.User, string, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	Close() error
}