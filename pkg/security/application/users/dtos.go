package users

type LoginCommand struct {
	UsernameOrEmail string `json:"username"`
	Password        string `json:"password"`
}

type RegisterCommand struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (cmd *RegisterCommand) Validate() error {
	return nil
}

type UpdateCommand struct {
	Name     string `json:"name"`
	Lastname string `json:"lastnaem"`
}

type ChangePasswordCommand struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
