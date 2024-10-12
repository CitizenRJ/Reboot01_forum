package user

type User struct {
	Uid      int    `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
