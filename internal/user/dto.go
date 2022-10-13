package user

type CreateUserDTO struct {
	Login          string `json:"login"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UpdateUserDTO struct {
	//Id                int   `json:"id"`
	//Login             string `json:"login"`
	OldPassword       string `json:"old_password"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
	//Password          string `json:"password"`
}

type AuthUserDTO struct {
	//Login    string `json:"login"`
	Password string `json:"password"`
}
