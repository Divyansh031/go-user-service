package handlers

import (
	"context"
	"log/slog"

	"github.com/Divyansh031/user-service/internal/domain"
	"github.com/Divyansh031/user-service/internal/storage"
	pb "github.com/Divyansh031/user-service/api/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	storage storage.Storage
}

func NewUserServiceServer(storage storage.Storage) *UserServiceServer {
	return &UserServiceServer{
		storage: storage,
	}
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	slog.Info("Creating user", "email", req.Email, "phone", req.PhoneNumber)

	user := domain.NewUser(
		req.FirstName,
		req.LastName,
		req.Gender,
		req.DateOfBirth.AsTime(),
		req.PhoneNumber,
		req.Email,
	)

	if err := user.Validate(); err != nil {
		slog.Error("Validation failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.storage.CreateUser(ctx, user); err != nil {
		if err == domain.ErrEmailAlreadyExists || err == domain.ErrPhoneAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		slog.Error("Failed to create user", "error", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	slog.Info("User created successfully", "user_id", user.ID)

	return &pb.CreateUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slog.Info("Getting user", "id", req.Id)

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		slog.Error("Failed to get user", "error", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// UpdateUser updates an existing user
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	slog.Info("Updating user", "id", req.Id)

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Only update fields that are provided (non-empty)
	firstName := req.FirstName
	if firstName == "" {
		firstName = user.FirstName
	}

	lastName := req.LastName
	if lastName == "" {
		lastName = user.LastName
	}

	gender := req.Gender
	if gender == "" {
		gender = user.Gender
	}

	dob := req.DateOfBirth.AsTime()
	if dob.IsZero() {
		dob = user.DateOfBirth
	}

	user.Update(firstName, lastName, gender, dob)

	if err := user.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		slog.Error("Failed to update user", "error", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	slog.Info("User updated successfully", "user_id", user.ID)

	return &pb.UpdateUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// DeleteUser deletes a user
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	slog.Info("Deleting user", "id", req.Id)

	if err := s.storage.DeleteUser(ctx, req.Id); err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		slog.Error("Failed to delete user", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	slog.Info("User deleted successfully", "user_id", req.Id)

	return &emptypb.Empty{}, nil
}

// BlockUser blocks a user
func (s *UserServiceServer) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	slog.Info("Blocking user", "id", req.Id)

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	if user.IsBlocked {
		return nil, status.Error(codes.FailedPrecondition, "user is already blocked")
	}

	user.Block()

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		slog.Error("Failed to block user", "error", err)
		return nil, status.Error(codes.Internal, "failed to block user")
	}

	slog.Info("User blocked successfully", "user_id", user.ID)

	return &pb.BlockUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// UnblockUser unblocks a user
func (s *UserServiceServer) UnblockUser(ctx context.Context, req *pb.UnblockUserRequest) (*pb.UnblockUserResponse, error) {
	slog.Info("Unblocking user", "id", req.Id)

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	if !user.IsBlocked {
		return nil, status.Error(codes.FailedPrecondition, "user is not blocked")
	}

	user.Unblock()

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		slog.Error("Failed to unblock user", "error", err)
		return nil, status.Error(codes.Internal, "failed to unblock user")
	}

	slog.Info("User unblocked successfully", "user_id", user.ID)

	return &pb.UnblockUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// UpdateUserContact updates user's phone number or email
func (s *UserServiceServer) UpdateUserContact(ctx context.Context, req *pb.UpdateUserContactRequest) (*pb.UpdateUserContactResponse, error) {
	slog.Info("Updating user contact", "id", req.Id)

	user, err := s.storage.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	var phone, email *string
	if req.PhoneNumber != nil {
		phone = req.PhoneNumber
	}
	if req.Email != nil {
		email = req.Email
	}

	user.UpdateContact(phone, email)

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		slog.Error("Failed to update user contact", "error", err)
		return nil, status.Error(codes.Internal, "failed to update user contact")
	}

	slog.Info("User contact updated successfully", "user_id", user.ID)

	return &pb.UpdateUserContactResponse{
		User: domainUserToProto(user),
	}, nil
}

// GetUserByPhone retrieves a user by phone number
func (s *UserServiceServer) GetUserByPhone(ctx context.Context, req *pb.GetUserByPhoneRequest) (*pb.GetUserResponse, error) {
	slog.Info("Getting user by phone", "phone", req.PhoneNumber)

	user, err := s.storage.GetUserByPhone(ctx, req.PhoneNumber)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		slog.Error("Failed to get user by phone", "error", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	slog.Info("Getting user by email", "email", req.Email)

	user, err := s.storage.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		slog.Error("Failed to get user by email", "error", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: domainUserToProto(user),
	}, nil
}

// ListUsers lists all users with pagination
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	slog.Info("Listing users", "page_size", req.PageSize)

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, nextToken, err := s.storage.ListUsers(ctx, pageSize, req.PageToken)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = domainUserToProto(user)
	}

	return &pb.ListUsersResponse{
		Users:         protoUsers,
		NextPageToken: nextToken,
		TotalCount:    int32(len(users)),
	}, nil
}

// Helper function to convert domain user to proto
func domainUserToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: timestamppb.New(user.DateOfBirth),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		IsBlocked:   user.IsBlocked,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}