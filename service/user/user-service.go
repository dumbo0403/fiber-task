package user

import (
	user_dto "github.com/dumbo0403/user_management/dto/user"
	user_entity "github.com/dumbo0403/user_management/entity/user"
	"github.com/dumbo0403/user_management/helper"
	user_repo "github.com/dumbo0403/user_management/repository/user"
	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByEmail(email string) (user_entity.User, error)
	IsDuplicatedEmail(email string) bool
	CreateUser(user user_entity.User) (user_entity.User, error)
	Login(user user_dto.UserLoginDTO) (user_dto.UserLoggedDTO, error)
	UpdateUser(user user_entity.User) (user_entity.User, error)
}

type userService struct {
	userRepo user_repo.UserRepository
}

func NewUserService(userRepo user_repo.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (service *userService) GetUserByEmail(email string) (user_entity.User, error) {
	return service.userRepo.GetUserByEmail(email)
}

func (service *userService) IsDuplicatedEmail(email string) bool {
	_, err := service.userRepo.GetUserByEmail(email)
	return err != gorm.ErrRecordNotFound
}

func (service *userService) CreateUser(user user_entity.User) (user_entity.User, error) {

	if service.IsDuplicatedEmail(user.Email) {
		return user_entity.User{}, gorm.ErrDuplicatedKey
	}

	return service.userRepo.CreateUser(user)
}

func (service *userService) Login(user user_dto.UserLoginDTO) (user_dto.UserLoggedDTO, error) {

	if service.IsDuplicatedEmail(user.Email) {
		return service.userRepo.Login(user)
	}
	return user_dto.UserLoggedDTO{}, gorm.ErrRecordNotFound
}

func (service *userService) UpdateUser(user user_entity.User) (user_entity.User, error) {
	updateUser := helper.BuildUpdateData(user)
	oldUser, err := service.userRepo.GetUserByID(user.ID)

	if err != nil {
		return user_entity.User{}, err
	}

	updateUser, err_update := service.userRepo.UpdateUser(updateUser)

	if err_update != nil {
		return user_entity.User{}, nil
	}

	err_mapping := smapping.FillStruct(&oldUser, updateUser)
	if err_mapping != nil {
		panic(err_mapping)
	}

	return oldUser, nil
}
