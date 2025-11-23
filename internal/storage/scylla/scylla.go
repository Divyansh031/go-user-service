// internal/storage/scylla/scylla.go
package scylla

import (
	"context"
	"fmt"
	"time"

	"github.com/Divyansh031/user-service/internal/domain"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type ScyllaDB struct {
	session *gocql.Session
}

func NewScyllaDB(hosts []string, port int, keyspace string, consistency string) (*ScyllaDB, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = parseConsistency(consistency)
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ScyllaDB: %w", err)
	}

	return &ScyllaDB{session: session}, nil
}

func parseConsistency(consistency string) gocql.Consistency {
	switch consistency {
	case "ONE":
		return gocql.One
	case "QUORUM":
		return gocql.Quorum
	case "ALL":
		return gocql.All
	default:
		return gocql.Quorum
	}
}

// CreateUser creates a new user
func (db *ScyllaDB) CreateUser(ctx context.Context, user *domain.User) error {
	exists, err := db.CheckEmailExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

	exists, err = db.CheckPhoneExists(ctx, user.PhoneNumber)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrPhoneAlreadyExists
	}

	query := `INSERT INTO users (id, first_name, last_name, gender, date_of_birth, 
		phone_number, email, is_blocked, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if err := db.session.Query(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Gender,
		user.DateOfBirth,
		user.PhoneNumber,
		user.Email,
		user.IsBlocked,
		user.CreatedAt,
		user.UpdatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	if err := db.insertPhoneLookup(ctx, user.PhoneNumber, user.ID); err != nil {
		return err
	}
	if err := db.insertEmailLookup(ctx, user.Email, user.ID); err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID (string)
func (db *ScyllaDB) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, first_name, last_name, gender, date_of_birth, 
		phone_number, email, is_blocked, created_at, updated_at 
		FROM users WHERE id = ?`

	var user domain.User
	if err := db.session.Query(query, id).WithContext(ctx).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.PhoneNumber,
		&user.Email,
		&user.IsBlocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByPhone retrieves a user by phone number
func (db *ScyllaDB) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var userID string
	query := `SELECT user_id FROM users_by_phone WHERE phone_number = ?`
	if err := db.session.Query(query, phone).WithContext(ctx).Scan(&userID); err != nil {
		if err == gocql.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to lookup user by phone: %w", err)
	}
	return db.GetUserByID(ctx, userID)
}

// GetUserByEmail retrieves a user by email
func (db *ScyllaDB) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userID string
	query := `SELECT user_id FROM users_by_email WHERE email = ?`
	if err := db.session.Query(query, email).WithContext(ctx).Scan(&userID); err != nil {
		if err == gocql.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to lookup user by email: %w", err)
	}
	return db.GetUserByID(ctx, userID)
}

// UpdateUser updates an existing user
func (db *ScyllaDB) UpdateUser(ctx context.Context, user *domain.User) error {
	existingUser, err := db.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}

	query := `UPDATE users SET first_name = ?, last_name = ?, gender = ?, 
		date_of_birth = ?, phone_number = ?, email = ?, is_blocked = ?, updated_at = ? 
		WHERE id = ?`

	if err := db.session.Query(query,
		user.FirstName, user.LastName, user.Gender,
		user.DateOfBirth, user.PhoneNumber, user.Email,
		user.IsBlocked, user.UpdatedAt, user.ID,
	).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if existingUser.PhoneNumber != user.PhoneNumber {
		db.deletePhoneLookup(ctx, existingUser.PhoneNumber)
		db.insertPhoneLookup(ctx, user.PhoneNumber, user.ID)
	}
	if existingUser.Email != user.Email {
		db.deleteEmailLookup(ctx, existingUser.Email)
		db.insertEmailLookup(ctx, user.Email, user.ID)
	}
	return nil
}

// DeleteUser deletes a user
func (db *ScyllaDB) DeleteUser(ctx context.Context, id string) error {
	user, err := db.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	db.deletePhoneLookup(ctx, user.PhoneNumber)
	db.deleteEmailLookup(ctx, user.Email)

	query := `DELETE FROM users WHERE id = ?`
	if err := db.session.Query(query, id).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers lists users with pagination
func (db *ScyllaDB) ListUsers(ctx context.Context, limit int, pageToken string) ([]*domain.User, string, error) {
	var users []*domain.User
	query := `SELECT id, first_name, last_name, gender, date_of_birth, 
		phone_number, email, is_blocked, created_at, updated_at FROM users LIMIT ?`

	iter := db.session.Query(query, limit).WithContext(ctx).Iter()
	var user domain.User
	for iter.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Gender,
		&user.DateOfBirth, &user.PhoneNumber, &user.Email, &user.IsBlocked,
		&user.CreatedAt, &user.UpdatedAt) {
		u := user
		users = append(users, &u)
	}
	if err := iter.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to list users: %w", err)
	}

	nextToken := ""
	if len(users) == limit {
		nextToken = users[len(users)-1].ID
	}
	return users, nextToken, nil
}

// CheckEmailExists checks if email exists
func (db *ScyllaDB) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var userID string
	query := `SELECT user_id FROM users_by_email WHERE email = ?`
	err := db.session.Query(query, email).WithContext(ctx).Scan(&userID)
	if err == gocql.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return true, nil
}

// CheckPhoneExists checks if phone exists
func (db *ScyllaDB) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	var userID string
	query := `SELECT user_id FROM users_by_phone WHERE phone_number = ?`
	err := db.session.Query(query, phone).WithContext(ctx).Scan(&userID)
	if err == gocql.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check phone: %w", err)
	}
	return true, nil
}

// Helper methods — now accept string userID
func (db *ScyllaDB) insertPhoneLookup(ctx context.Context, phone string, userID string) error {
	uid := uuid.MustParse(userID)
	query := `INSERT INTO users_by_phone (phone_number, user_id, created_at) VALUES (?, ?, ?)`
	return db.session.Query(query, phone, uid, time.Now()).WithContext(ctx).Exec()
}

func (db *ScyllaDB) deletePhoneLookup(ctx context.Context, phone string) error {
	query := `DELETE FROM users_by_phone WHERE phone_number = ?`
	return db.session.Query(query, phone).WithContext(ctx).Exec()
}

func (db *ScyllaDB) insertEmailLookup(ctx context.Context, email string, userID string) error {
	uid := uuid.MustParse(userID)
	query := `INSERT INTO users_by_email (email, user_id, created_at) VALUES (?, ?, ?)`
	return db.session.Query(query, email, uid, time.Now()).WithContext(ctx).Exec()
}

func (db *ScyllaDB) deleteEmailLookup(ctx context.Context, email string) error {
	query := `DELETE FROM users_by_email WHERE email = ?`
	return db.session.Query(query, email).WithContext(ctx).Exec()
}

func (db *ScyllaDB) Close() error {
	db.session.Close()
	return nil
}