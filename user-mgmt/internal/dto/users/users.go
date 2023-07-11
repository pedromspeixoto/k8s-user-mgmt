package users

import (
	usermodel "github.com/pedromspeixoto/users-api/internal/data/models/users"
	"github.com/pedromspeixoto/users-api/internal/pkg/uuid"
)

// request
type UserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func ModelFromUserRequest(post *UserRequest) *usermodel.User {
	model := &usermodel.User{
		UserId: uuid.GenerateUUID(),
		Email:  post.Email,
	}
	return model
}

// response
type UserResponse struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

func NewUserResponse(user *usermodel.User) *UserResponse {
	resp := &UserResponse{
		UserId: user.UserId,
		Email:  user.Email,
	}
	return resp
}

type UserListResponse struct {
	Users []UserResponse `json:"users,omitempty"`
}

func NewUserListResponse(models []usermodel.User) *UserListResponse {
	var users []UserResponse
	for _, m := range models {
		users = append(users, UserResponse{
			UserId: m.UserId,
			Email:  m.Email,
		})
	}
	return &UserListResponse{Users: users}
}
