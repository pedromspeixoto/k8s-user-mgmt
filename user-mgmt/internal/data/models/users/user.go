package users

import (
	"math"

	"github.com/pedromspeixoto/users-api/internal/data"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId string
	Email  string
}

// UserRepository is a repository for dealing with the user object.
type UserRepository interface {
	// List lists users from the database with pagination.
	List(limit, page int) ([]User, *data.Pagination, error)
	// GetByUUID gets a user from the database by uuid.
	GetByUUID(uuid string) (*User, error)
	// Get gets a user from the database by id.
	Get(id uint) (*User, error)
	// Create creates a user in the database.
	Create(user *User) error
	// SoftDelete soft deletes a user record from the database.
	SoftDelete(user *User) error
	// HardDelete hard deletes a user record from the database.
	HardDelete(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) List(limit, page int) ([]User, *data.Pagination, error) {
	var users []User

	// pagination object
	pagination := &data.Pagination{
		Limit: limit,
		Page:  page,
	}
	u.db.Scopes(pagination.Paginate()).Find(&users)

	// pagination details
	u.db.Model(&User{}).Count(&pagination.TotalRows)
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))

	return users, pagination, nil
}

func (u userRepository) GetByUUID(uuid string) (*User, error) {
	user := User{}
	result := u.db.Unscoped().Where("user_id = ?", uuid).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (u userRepository) Get(id uint) (*User, error) {
	user := User{}
	result := u.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u userRepository) Create(user *User) error {
	result := u.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u userRepository) SoftDelete(user *User) error {
	result := u.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u userRepository) HardDelete(user *User) error {
	result := u.db.Unscoped().Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
