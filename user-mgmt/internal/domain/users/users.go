package users

import (
	"context"
	"fmt"
	"github.com/pedromspeixoto/users-api/internal/data/models/users"
	"github.com/pedromspeixoto/users-api/internal/pkg/files"
	"github.com/pedromspeixoto/users-api/internal/pkg/uuid"
	"net/http"
	"strings"

	"github.com/pedromspeixoto/users-api/internal/config"
	"github.com/pedromspeixoto/users-api/internal/dto"
	usersdto "github.com/pedromspeixoto/users-api/internal/dto/users"
	"github.com/pedromspeixoto/users-api/internal/pkg/logger"
	"go.uber.org/fx"
)

const (
	PdfFileType = "application/pdf"
)

// UserService provides methods pertaining to managing users.
type UserService interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, Post *usersdto.UserRequest) (int, *usersdto.UserResponse, error)
	// ListUsers retrieves all users with pagination.
	ListUsers(ctx context.Context, pagination *dto.PaginationRequest) (int, *dto.PaginationResponse, error)
	// GetUser retrieves a user by uuid
	GetUser(ctx context.Context, uuid string) (int, *usersdto.UserResponse, error)
	// DeleteUser soft deletes a user entry by uuid
	DeleteUser(ctx context.Context, uuid string) (int, error)

	// CreateUserFile creates a new user file
	CreateUserFile(ctx context.Context, userId string) (int, *usersdto.UserFileResponse, error)
	// ListUserFiles retrieves all user files with pagination.
	ListUserFiles(ctx context.Context, userId string, pagination *dto.PaginationRequest) (int, *dto.PaginationResponse, error)
	// GetUserFile retrieves a user file by uuid
	GetUserFile(ctx context.Context, userId string, uuid string) (int, *usersdto.UserFileResponse, error)
	// DeleteUserFile soft deletes a user file entry by uuid
	DeleteUserFile(ctx context.Context, userId string, uuid string) (int, error)
	// DownloadUserFile download a user file by uuid
	DownloadUserFile(ctx context.Context, userId string, uuid string) (int, []byte, string, error)
}

type UserServiceDeps struct {
	fx.In

	Config             *config.Config
	Logger             *logger.LoggingClient
	UserRepository     users.UserRepository
	UserFileRepository users.UserFileRepository
	FileServingClient  *files.FileServingClient
}

type userService struct {
	UserServiceDeps
	logger.Logger
}

func NewUserService(deps UserServiceDeps) UserService {
	return &userService{
		UserServiceDeps: deps,
		Logger:          deps.Logger.GetLogger(),
	}
}

func (u *userService) CreateUser(ctx context.Context, request *usersdto.UserRequest) (int, *usersdto.UserResponse, error) {
	model := usersdto.ModelFromUserRequest(request)
	err := u.UserRepository.Create(model)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error creating new user: %v", err)
	}

	return http.StatusCreated, usersdto.NewUserResponse(model), nil
}

func (u *userService) ListUsers(ctx context.Context, paginationRequest *dto.PaginationRequest) (int, *dto.PaginationResponse, error) {
	users, pageEnv, err := u.UserRepository.List(paginationRequest.Limit, paginationRequest.Page)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpdected error fetching users: %v", err)
	}

	pageEnv.Data = usersdto.NewUserListResponse(users)
	return http.StatusOK, dto.NewPaginationResponse(pageEnv), nil
}

func (u *userService) GetUser(ctx context.Context, uuid string) (int, *usersdto.UserResponse, error) {
	user, err := u.UserRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, usersdto.NewUserResponse(user), nil
}

func (u *userService) DeleteUser(ctx context.Context, uuid string) (int, error) {
	post, err := u.UserRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, err
	}

	err = u.UserRepository.SoftDelete(post)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unexpected error deleting user: %v", err)
	}

	return http.StatusOK, nil
}

func (u *userService) CreateUserFile(ctx context.Context, userId string) (int, *usersdto.UserFileResponse, error) {
	// check if user exists
	code, user, err := u.GetUser(ctx, userId)
	if err != nil {
		return code, nil, err
	}

	fileContent, fileType, err := u.FileServingClient.GetRandomFile()
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error fetching file from client: %v", err)
	}

	if strings.Contains(fileType, PdfFileType) {
		err = files.CheckPDFCorrupted(fileContent)
		if err != nil {
			return http.StatusBadRequest, nil, fmt.Errorf("file possibly corrupted. could not open file: %v", err)
		}
	}

	model := &users.UserFile{
		UserId:      user.UserId,
		FileId:      uuid.GenerateUUID(),
		FileType:    fileType,
		FileContent: fileContent,
	}

	err = u.UserFileRepository.Create(model)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error creating new user file: %v", err)
	}

	return http.StatusCreated, usersdto.NewUserFileResponse(model), nil
}

func (u *userService) ListUserFiles(ctx context.Context, userId string, paginationRequest *dto.PaginationRequest) (int, *dto.PaginationResponse, error) {
	// check if user exists
	code, _, err := u.GetUser(ctx, userId)
	if err != nil {
		return code, nil, err
	}

	files, pageEnv, err := u.UserFileRepository.List(userId, paginationRequest.Limit, paginationRequest.Page)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpdected error fetching user files: %v", err)
	}

	pageEnv.Data = usersdto.NewUserFileListResponse(files)
	return http.StatusOK, dto.NewPaginationResponse(pageEnv), nil
}

func (u *userService) GetUserFile(ctx context.Context, userId string, fileId string) (int, *usersdto.UserFileResponse, error) {
	// check if user exists
	code, _, err := u.GetUser(ctx, userId)
	if err != nil {
		return code, nil, err
	}

	file, err := u.UserFileRepository.GetByUUID(fileId)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, usersdto.NewUserFileResponse(file), nil
}

func (u *userService) DeleteUserFile(ctx context.Context, userId string, uuid string) (int, error) {
	// check if user exists
	code, _, err := u.GetUser(ctx, userId)
	if err != nil {
		return code, err
	}

	file, err := u.UserFileRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, err
	}

	err = u.UserFileRepository.SoftDelete(file)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unexpected error deleting user file: %v", err)
	}

	return http.StatusOK, nil
}

func (u *userService) DownloadUserFile(ctx context.Context, userId string, uuid string) (int, []byte, string, error) {
	// check if user exists
	code, _, err := u.GetUser(ctx, userId)
	if err != nil {
		return code, nil, "", err
	}

	file, err := u.UserFileRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, nil, "", err
	}

	// get file extension from file type
	fileExtension := strings.Split(file.FileType, "/")[1]

	return http.StatusOK, file.FileContent, fileExtension, nil
}
