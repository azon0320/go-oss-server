package permission

import "github.com/dormao/go-oss-server/internal/utils"

type PermMap struct {
	OwnPermissions []string
	AllPermission  bool
}

func (m *PermMap) HasPermissions(need []string) bool {
	if m.AllPermission {
		return true
	}
	var count = 0
	for _, needPerm := range need {
		if utils.InArrayS(needPerm, m.OwnPermissions) {
			count++
		}
	}
	return count >= len(need)
}

func (m *PermMap) AddPermissions(perms []string) {
	if !m.AllPermission {
		for _, perm := range perms {
			if !utils.InArrayS(perm, m.OwnPermissions) {
				m.OwnPermissions = append(m.OwnPermissions, perm)
			}
		}
	}
}

func (m *PermMap) PermissionOutput() []string {
	if m.AllPermission {
		return []string{"*"}
	} else {
		return m.OwnPermissions
	}
}

func FromStringArray(perms []string) *PermMap {
	var AllPerm = utils.InArrayS("*", perms)
	var themap = &PermMap{
		OwnPermissions: []string{},
		AllPermission:  AllPerm,
	}
	if !AllPerm {
		themap.OwnPermissions = perms
	}
	return themap
}
