package permission

import "testing"

func TestSU(t *testing.T) {
	su := CreateSuperUser()
	t.Logf("su has perm: %t", su.HasPermissions([]string{
		"user.create", "user.delete",
	}))
	u := FromStringArray([]string{
		"user.create",
	})
	t.Logf("user has perm: %t", u.HasPermissions([]string{
		"user.create",
	}))
}
