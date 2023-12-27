package main

import (
	"github.com/dumbo0403/user_management/cache"
	"github.com/dumbo0403/user_management/config"
	user_controller "github.com/dumbo0403/user_management/controller/user"
	jwtAuth "github.com/dumbo0403/user_management/jwt"
	user_repo "github.com/dumbo0403/user_management/repository/user"
	user_service "github.com/dumbo0403/user_management/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                 = config.SetupDataBaseConnection()
	rdb            *redis.Client            = cache.ConnectionRedisDB()
	UserRepository user_repo.UserRepository = user_repo.NewUserRepository(db, rdb)

	UserService user_service.UserService = user_service.NewUserService(UserRepository)

	UserController user_controller.UserController = user_controller.NewUserController(db, UserService, rdb)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	defer cache.CloseRedis(rdb)
	app := fiber.New()

	app.Post("/login", UserController.Login)
	app.Post("/create_random_user", UserController.CreateDemoUser)
	app.Use(jwtAuth.ValidateToken())

	app.Get("/", UserController.GetAllUsers)
	app.Post("/insert", UserController.SignUp)
	app.Patch("/update", UserController.UpdateUser)
	app.Delete("/delete", UserController.DeleteUser)

	app.Listen(":3000")
}
