package roles

type Module struct {
	Code string
	Name string
}

type Permission struct {
	Permission string
	Module     Module
}

type Role struct {
	Code        string
	Name        string
	Permissions []Permission
}
