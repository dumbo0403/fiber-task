package user

import (
	"fmt"
	"strconv"

	"github.com/dumbo0403/user_management/cache"
	user_dto "github.com/dumbo0403/user_management/dto/user"
	user_entity "github.com/dumbo0403/user_management/entity/user"
	user_service "github.com/dumbo0403/user_management/service/user"
	"github.com/redis/go-redis/v9"

	jwtAuth "github.com/dumbo0403/user_management/jwt"

	"github.com/dumbo0403/user_management/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController interface {
	GetAllUsers(c *fiber.Ctx) error
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	CreateDemoUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userController struct {
	db          *gorm.DB
	rdb         *redis.Client
	userService user_service.UserService
}

func NewUserController(db *gorm.DB, userService user_service.UserService, rdb *redis.Client) UserController {
	return &userController{
		db:          db,
		rdb:         rdb,
		userService: userService,
	}
}

func (ctrl *userController) GetAllUsers(c *fiber.Ctx) error {

	users := []user_entity.User{}
	err := ctrl.db.Find(&users).Error
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))

	}
	return c.JSON(helper.BuildResponse(users))
}

func (ctrl *userController) SignUp(c *fiber.Ctx) error {
	user := user_entity.User{}
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))
	}

	user, err_create := ctrl.userService.CreateUser(user)

	if err_create != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err_create.Error()))
	}
	return c.JSON(helper.BuildResponse(user))
}

func (ctrl *userController) Login(c *fiber.Ctx) error {
	user := user_dto.UserLoginDTO{}
	var user_login user_dto.UserLoggedDTO
	err_parse := c.BodyParser(&user)

	if err_parse != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err_parse.Error()))
	}

	user_login, err_redis := cache.GetUser(ctrl.rdb, user)

	if err_redis != nil {
		user_login, err_login := ctrl.userService.Login(user)

		if err_login != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err_login))
		}

		_, err_set := cache.SetUser(ctrl.rdb, user.Email, user_login)

		if err_set != nil {
			fmt.Println("err on redis", err_set)
		}

		fmt.Println("User successfully inserted into redis", user_login)

		token, err := jwtAuth.GenerateJWT(user_login.FirstName, user_login.ID)

		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))
		}

		user_login.Token = token

		return c.JSON(helper.BuildResponse(user_login))
	}

	token, err := jwtAuth.GenerateJWT(user_login.FirstName, user_login.ID)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))
	}

	user_login.Token = token

	return c.JSON(helper.BuildResponse(user_login))
}

func (ctrl *userController) UpdateUser(c *fiber.Ctx) error {

	user := user_entity.User{}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))
	}

	updatedUser, err_update := ctrl.userService.UpdateUser(user)

	if err_update != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err_update.Error()))
	}

	return c.JSON(helper.BuildResponse(updatedUser))
}

func (ctrl *userController) DeleteUser(c *fiber.Ctx) error {

	user_id, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	deletedUser := user_entity.User{}

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err.Error()))
	}
	ctrl.db.Take(&deletedUser, "id = ?", user_id)

	cache.DeleteUser(ctrl.rdb, deletedUser.Email)
	err_delete := ctrl.db.Delete(&deletedUser, "id = ?", user_id).Error

	if err_delete != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(helper.BuildErrorResponse(err_delete.Error()))
	}

	return c.JSON(helper.BuildResponse(deletedUser))
}

func (ctrl *userController) CreateDemoUser(c *fiber.Ctx) error {
	demoUser := user_entity.User{}
	demoUser.Email = "demo@mail.com"
	demoUser.FirstName = "demo"
	demoUser.LastName = "demo"
	demoUser.Password = "demo"

	user, _ := ctrl.userService.CreateUser(demoUser)
	return c.JSON(helper.BuildErrorResponse(user))
}
