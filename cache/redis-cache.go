package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	user_dto "github.com/dumbo0403/user_management/dto/user"
)

var ctx context.Context = context.Background()

func ConnectionRedisDB() *redis.Client {
	// ctx = context.Background()
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		},
	)

	result, err := rdb.Ping(ctx).Result()

	fmt.Println(err, result)

	return rdb
}

func CloseRedis(rdb *redis.Client) {
	err := rdb.Close()
	if err != nil {
		panic(err)
	}
}

func SetUser(rdb *redis.Client, email string, user user_dto.UserLoggedDTO) (user_dto.UserLoggedDTO, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return user_dto.UserLoggedDTO{}, err
	}

	_, errSet := rdb.Set(ctx, email, userJson, time.Hour).Result()

	if errSet != nil {
		return user_dto.UserLoggedDTO{}, errSet
	}

	return user, nil
}

func GetUser(rdb *redis.Client, logged_user user_dto.UserLoginDTO) (user_dto.UserLoggedDTO, error) {
	user := user_dto.UserLoggedDTO{}
	result, err := rdb.Get(ctx, logged_user.Email).Result()

	if err != nil {
		return user_dto.UserLoggedDTO{}, err
	}
	errUnMarshal := json.Unmarshal([]byte(result), &user)

	if errUnMarshal != nil {
		return user_dto.UserLoggedDTO{}, errUnMarshal
	}

	if !comparePassword(user.Password, []byte(logged_user.Password)) {
		return user_dto.UserLoggedDTO{}, errors.New("WrongPassword")
	}

	return user, nil
}

func UpdateEmail(rdb *redis.Client, oldEmail, newEmail string) error {
	val, err := rdb.Get(ctx, oldEmail).Result()
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, newEmail, val, time.Hour).Err()
	if err != nil {
		return err
	}

	err = rdb.Del(ctx, oldEmail).Err()
	return err
}

func DeleteUser(rdb *redis.Client, email string) error {
	_, err := rdb.Del(ctx, email).Result()

	if err != nil {
		return err
	}

	return nil
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
