package user

type User struct {
	ID        uint64 ` json:"id"`
	FirstName string ` json:"first_name"`
	LastName  string ` json:"last_name"`
	Email     string ` json:"email"`
	Password  string ` json:"password"`
}

func (User) TableName() string {
	return "user"
}
