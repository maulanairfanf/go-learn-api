package repositories

import (
	"myapi/db"
	"myapi/models"
	"myapi/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct{}

func (r *UserRepository) FindAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.User, int64, error) {
	var user []models.User
	var total int64

	query := db.DB.Model(&models.User{})

	if search != "" {
		query = query.Where("username ILIKE ?", "%"+search+"%")
	}

	if sort != "" {
		query = query.Order(sort + " " + order)
	}

	query = utils.ApplyPagination(query, pagination, &total)

	err := query.Find(&user).Error

	if !pagination.IsPaginated() {
		total = int64(len(user))
	}

	return user, total, err
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepository) Save(user *models.User) error {
	return db.DB.Save(user).Error
}

func (r *UserRepository) Delete(user *models.User) error {
	return db.DB.Delete(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) MatchPassword(password, hash string) bool {
	return checkPasswordHash(password, hash)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
