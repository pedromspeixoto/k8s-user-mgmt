package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pedromspeixoto/users-api/internal/config"
	"github.com/pedromspeixoto/users-api/internal/domain/users"
	"github.com/pedromspeixoto/users-api/internal/dto"
	usersdto "github.com/pedromspeixoto/users-api/internal/dto/users"
	"github.com/pedromspeixoto/users-api/internal/http/handlers/common"
	"github.com/pedromspeixoto/users-api/internal/http/middlewares"
	"github.com/pedromspeixoto/users-api/internal/pkg/logger"
	"go.uber.org/fx"
)

type UserServiceHandler interface {
	Routes() chi.Router
}

type userServiceDeps struct {
	fx.In

	Config      *config.Config
	Logger      *logger.LoggingClient
	Validator   *validator.Validate
	UserService users.UserService
}

type userServiceHandler struct {
	userServiceDeps
	logger.Logger
}

func NewUserServiceHandler(deps userServiceDeps) UserServiceHandler {
	return &userServiceHandler{
		userServiceDeps: deps,
		Logger:          deps.Logger.GetLogger(),
	}
}

func (h userServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// users
	r.With(middlewares.Paginate).Get("/", h.ListUsers)
	r.Post("/", h.CreateUser)
	r.Put("/", h.CreateUser)
	r.Get("/{userId}", h.GetUser)
	r.Delete("/{userId}", h.DeleteUser)

	// user files
	r.With(middlewares.Paginate).Get("/{userId}/files", h.ListUserFiles)
	r.Post("/{userId}/files", h.CreateUserFile)
	r.Get("/{userId}/files/{fileId}", h.GetUserFile)
	r.Delete("/{userId}/files/{fileId}", h.DeleteUserFile)
	r.Get("/{userId}/files/{fileId}/download", h.DownloadUserFile)

	return r
}

// ListUsers - Handles user management
// @Summary Gets all users.
// @Description This API is used to list all users
// @Param limit query int false "Limit"
// @Param page  query int false "Page"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users [get]
func (h userServiceHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(middlewares.LimitKey).(int)
	page := r.Context().Value(middlewares.PageKey).(int)
	sort := r.Context().Value(middlewares.SortKey).(string)
	filter := r.Context().Value(middlewares.FilterKey).(map[string]string)
	search := r.Context().Value(middlewares.SearchKey).(map[string]string)

	pageRequest, err := dto.NewPaginationRequest(limit, page, sort, filter, search)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
	}

	statusCode, env, err := h.userServiceDeps.UserService.ListUsers(r.Context(), pageRequest)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "users retrieved", env)
}

// CreateUser - Handles user management
// @Summary Create a new user.
// @Description This API is used to create a new user
// @Param request body usersdto.UserRequest true "User Payload"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users [post]
func (h userServiceHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := usersdto.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		common.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Validator.Struct(user)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, userResponse, err := h.UserService.CreateUser(r.Context(), &user)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "new user created", userResponse)
}

// GetUser - Handles user management
// @Summary Get a user.
// @Description This API is used to get users by id
// @Param user_id path string true "User ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id} [get]
func (h userServiceHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	statusCode, user, err := h.userServiceDeps.UserService.GetUser(r.Context(), userId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "user retrieved", user)
}

// DeleteUser - Handles users mgmt
// @Summary Delete a user.
// @Description This API is used to delete a user
// @Param user_id path string true "User ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id} [delete]
func (h userServiceHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	statusCode, err := h.userServiceDeps.UserService.DeleteUser(r.Context(), userId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "", nil)
}

// ListUserFiles - Handles user files management
// @Summary Gets all user files.
// @Description This API is used to list all user files
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit"
// @Param page  query int false "Page"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id}/files [get]
func (h userServiceHandler) ListUserFiles(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	limit := r.Context().Value(middlewares.LimitKey).(int)
	page := r.Context().Value(middlewares.PageKey).(int)
	sort := r.Context().Value(middlewares.SortKey).(string)
	filter := r.Context().Value(middlewares.FilterKey).(map[string]string)
	search := r.Context().Value(middlewares.SearchKey).(map[string]string)

	pageRequest, err := dto.NewPaginationRequest(limit, page, sort, filter, search)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
	}

	statusCode, userFiles, err := h.userServiceDeps.UserService.ListUserFiles(r.Context(), userId, pageRequest)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "user files retrieved", userFiles)
}

// CreateUserFile - Handles user files management
// @Summary Create a new user file.
// @Description This API is used to create a new user file
// @Param user_id path string true "User ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id}/files [post]
func (h userServiceHandler) CreateUserFile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	statusCode, userFile, err := h.userServiceDeps.UserService.CreateUserFile(r.Context(), userId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "new user file created", userFile)
}

// GetUserFile - Handles user files management
// @Summary Get a user file.
// @Description This API is used to get user files by id
// @Param user_id path string true "User ID"
// @Param file_id path string true "File ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id}/files/{file_id} [get]
func (h userServiceHandler) GetUserFile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	fileId := chi.URLParam(r, "fileId")

	statusCode, userFile, err := h.userServiceDeps.UserService.GetUserFile(r.Context(), userId, fileId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "user file retrieved", userFile)
}

// DeleteUserFile - Handles user files management
// @Summary Delete a user file.
// @Description This API is used to delete a user file
// @Param user_id path string true "User ID"
// @Param file_id path string true "File ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id}/files/{file_id} [delete]
func (h userServiceHandler) DeleteUserFile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	fileId := chi.URLParam(r, "fileId")

	statusCode, err := h.userServiceDeps.UserService.DeleteUserFile(r.Context(), userId, fileId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "", nil)
}

// DownloadUserFile - Handles user files management
// @Summary Download a user file.
// @Description This API is used to download a user file
// @Param user_id path string true "User ID"
// @Param file_id path string true "File ID"
// @Tags users
// @Accept  json
// @Produce  json
// @Router /v1/users/{user_id}/files/{file_id}/download [get]
func (h userServiceHandler) DownloadUserFile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	fileId := chi.URLParam(r, "fileId")

	statusCode, file, fileType, err := h.userServiceDeps.UserService.DownloadUserFile(r.Context(), userId, fileId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	// set headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=filename.%s", fileType))
	w.Header().Set("Content-Type", http.DetectContentType(file))

	_, err = w.Write(file)
	if err != nil {
		common.Err(w, http.StatusInternalServerError, err.Error())
		return
	}
}
