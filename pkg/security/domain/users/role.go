package users

const (
	ADMIN string = "admin"
	USER         = "user"
)

const (
	CREATE string = "C"
	READ          = "R"
	UPDATE        = "U"
	DELETE        = "D"
)

type Permission struct {
	Permission string
	Module     string
}

type Role struct {
	Code        string
	Permissions []Permission
}
