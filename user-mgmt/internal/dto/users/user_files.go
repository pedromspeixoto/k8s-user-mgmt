package users

import (
	usermodel "github.com/pedromspeixoto/users-api/internal/data/models/users"
	"time"
)

// response
type UserFileResponse struct {
	FileId    string    `json:"file_id"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserFileResponse(userFile *usermodel.UserFile) *UserFileResponse {
	resp := &UserFileResponse{
		FileId:    userFile.FileId,
		FileType:  userFile.FileType,
		CreatedAt: userFile.CreatedAt,
	}
	return resp
}

type UserFileListResponse struct {
	UserFiles []UserFileResponse `json:"user_files,omitempty"`
}

func NewUserFileListResponse(models []usermodel.UserFile) *UserFileListResponse {
	var userFiles []UserFileResponse
	for _, m := range models {
		userFiles = append(userFiles, UserFileResponse{
			FileId:    m.FileId,
			FileType:  m.FileType,
			CreatedAt: m.CreatedAt,
		})
	}
	return &UserFileListResponse{UserFiles: userFiles}
}
