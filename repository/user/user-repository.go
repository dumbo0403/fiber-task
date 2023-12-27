package user

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/dumbo0403/user_management/cache"
	user_dto "github.com/dumbo0403/user_management/dto/user"
	user_entity "github.com/dumbo0403/user_management/entity/user"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint64) (user_entity.User, error)
	GetUserByEmail(email string) (user_entity.User, error)
	CreateUser(user user_entity.User) (user_entity.User, error)
	Login(user user_dto.UserLoginDTO) (user_dto.UserLoggedDTO, error)
	UpdateUser(updateUser map[string]interface{}) (map[string]interface{}, error)
}

type userRepository struct {
	connection *gorm.DB
	rdb        *redis.Client
}

func NewUserRepository(connection *gorm.DB, rdb *redis.Client) UserRepository {
	return &userRepository{
		connection: connection,
		rdb:        rdb,
	}
}

func (repo *userRepository) GetUserByEmail(email string) (user_entity.User, error) {
	user := user_entity.User{}

	err := repo.connection.Take(&user, "email = ?", email).Error

	if err != nil {
		return user_entity.User{}, err
	}

	return user, nil
}

func (repo *userRepository) CreateUser(user user_entity.User) (user_entity.User, error) {
	user.Password = hashAndSaltUser([]byte(user.Password))
	err_create := repo.connection.Create(&user).Error
	if err_create != nil {
		return user_entity.User{}, err_create
	}

	return user, nil
}

func (repo *userRepository) Login(login_user user_dto.UserLoginDTO) (user_dto.UserLoggedDTO, error) {
	user := user_dto.UserLoggedDTO{}

	err := repo.connection.Take(&user, "email = ?", login_user.Email).Error
	if err != nil {
		return user_dto.UserLoggedDTO{}, err
	}

	if comparePassword(user.Password, []byte(login_user.Password)) {
		return user, nil
	}

	return user_dto.UserLoggedDTO{}, errors.New("WrongPassword")
}

func (repo *userRepository) GetUserByID(id uint64) (user_entity.User, error) {
	user := user_entity.User{}

	if err := repo.connection.Take(&user, "id = ?", id).Error; err != nil {
		return user_entity.User{}, err
	}
	return user, nil
}

func (repo *userRepository) UpdateUser(updateUser map[string]interface{}) (map[string]interface{}, error) {

	if updateUser["Email"] != nil {
		oldUser := user_dto.UserLoggedDTO{}
		repo.connection.Take(&oldUser, "id = ?", updateUser["ID"])

		err := cache.UpdateEmail(repo.rdb, oldUser.Email, updateUser["Email"].(string))

		if err != nil {
			return nil, err
		}
	}

	if updateUser["Password"] != nil {
		updateUser["Password"] = hashAndSaltUser([]byte(updateUser["Password"].(string)))
	}

	err_update := repo.connection.Model(&user_entity.User{}).Where("id = ? ", updateUser["ID"]).Updates(&updateUser).Error

	if err_update != nil {
		return nil, err_update
	}

	updated_user, err := mapToStruct(updateUser)

	if err != nil {
		return nil, err
	}

	cache.SetUser(repo.rdb, updateUser["Email"].(string), updated_user)

	return updateUser, nil
}

func hashAndSaltUser(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func mapToStruct(m map[string]interface{}) (user_dto.UserLoggedDTO, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return user_dto.UserLoggedDTO{}, err
	}

	var result user_dto.UserLoggedDTO
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return user_dto.UserLoggedDTO{}, err
	}

	return result, nil
}
