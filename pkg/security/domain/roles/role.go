package roles

import "strings"

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
	permission string
	module     string
}

type Role struct {
	code        string
	permissions []Permission
}

func (r *Role) Code() string {
	return r.code
}

func (r *Role) HasPermissions(permissions string, module string) bool {
	for _, rolePerm := range r.permissions {
		if rolePerm.module == module {
			for _, perm := range strings.Split(permissions, "") {
				if !strings.Contains(rolePerm.permission, perm) {
					return false
				}
			}
			return true
		}
	}
	return false
}
