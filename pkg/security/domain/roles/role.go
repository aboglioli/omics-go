package roles

import "strings"

const (
	ADMIN           string = "admin"
	CONTENT_MANAGER        = "content-manager"
	USER                   = "user"
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

func (r *Role) Equals(role interface{}) bool {
	if ro, ok := role.(*Role); ok {
		return r.code == ro.code
	} else if co, ok := role.(string); ok {
		return r.code == co
	}
	return false
}

func (r *Role) Code() string {
	return r.code
}

func (r *Role) Is(code string) bool {
	return r.code == code
}

func (r *Role) HasPermissions(permissions string, module string) bool {
	if r.Is(ADMIN) {
		return true
	}

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
