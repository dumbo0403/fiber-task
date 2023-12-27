package user

type UserLoginDTO struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type UserLoggedDTO struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Password  string `json:"password"`
}

func (UserLoggedDTO) TableName() string {
	return "user"
}
