package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	subRole := New("sub_role")
	table := []struct {
		role *Role
	}{
		{
			role: New("test_role"),
		},
	}

	for _, tt := range table {
		err := tt.role.AddRole(subRole)
		assert.NoError(t, err)

		eff := tt.role.EffectiveRoles()
		assert.Greater(t, len(eff), 1)

		err = tt.role.AddRole(tt.role)
		assert.Error(t, err, "should not be able to add subrole that matches parent")
		err = tt.role.AddRole(subRole)
		assert.Error(t, err, "should not be able to add duplicate subrole")

		b, err := Match(tt.role, "test*")
		assert.NoError(t, err)
		assert.True(t, b)
	}
}
