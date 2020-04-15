package permission

func CreateSuperUser() *PermMap {
	return FromStringArray([]string{"*"})
}
