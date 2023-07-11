package users

import (
	"math"

	"github.com/pedromspeixoto/users-api/internal/data"
	"gorm.io/gorm"
)

type UserFile struct {
	gorm.Model
	UserId      string
	FileId      string
	FileType    string
	FileContent []byte
}

// UserFileRepository is a repository for dealing with user files.
type UserFileRepository interface {
	// List user files from the database with pagination.
	List(userId string, limit, page int) ([]UserFile, *data.Pagination, error)
	// GetByUserUUID gets a file from the database by user uuid.
	GetByUserUUID(uuid string) (*UserFile, error)
	// GetByUUID gets a file from the database by uuid.
	GetByUUID(uuid string) (*UserFile, error)
	// Get gets a file from the database by id.
	Get(id uint) (*UserFile, error)
	// Create creates a file in the database.
	Create(file *UserFile) error
	// SoftDelete soft deletes a file from the database.
	SoftDelete(file *UserFile) error
	// HardDelete hard deletes a file from the database.
	HardDelete(file *UserFile) error
}

type userFileRepository struct {
	db *gorm.DB
}

func NewUserFileRepository(db *gorm.DB) UserFileRepository {
	return &userFileRepository{
		db: db,
	}
}

func (f userFileRepository) List(userId string, limit, page int) ([]UserFile, *data.Pagination, error) {
	var userFiles []UserFile

	// pagination object
	pagination := &data.Pagination{
		Limit: limit,
		Page:  page,
	}
	f.db.Scopes(pagination.Paginate()).Find(&userFiles).Where("user_id = ?", userId)

	// pagination details
	f.db.Model(&UserFile{}).Count(&pagination.TotalRows)
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))

	return userFiles, pagination, nil
}

func (f userFileRepository) GetByUserUUID(uuid string) (*UserFile, error) {
	userFile := UserFile{}
	result := f.db.Unscoped().Where("user_id = ?", uuid).Find(&userFile)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &userFile, nil
}

func (f userFileRepository) GetByUUID(uuid string) (*UserFile, error) {
	userFile := UserFile{}
	result := f.db.Unscoped().Where("file_id = ?", uuid).Find(&userFile)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &userFile, nil
}

func (f userFileRepository) Get(id uint) (*UserFile, error) {
	userFile := UserFile{}
	result := f.db.First(&userFile, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userFile, nil
}

func (f userFileRepository) Create(userFile *UserFile) error {
	result := f.db.Create(userFile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (f userFileRepository) SoftDelete(userFile *UserFile) error {
	result := f.db.Delete(userFile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (f userFileRepository) HardDelete(userFile *UserFile) error {
	result := f.db.Unscoped().Delete(userFile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
