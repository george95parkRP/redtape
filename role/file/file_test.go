package file_test

import (
	"os"
	"testing"

	"github.com/blushft/redtape/role"
	"github.com/blushft/redtape/role/file"
	"github.com/stretchr/testify/assert"
)

func TestFileRoleManager(t *testing.T) {
	rm, err := file.New()
	if err != nil {
		t.Fatal(err)
	}

	r := role.New("test_role")
	if err := rm.Create(r); err != nil {
		t.Fatal(err)
	}

	r.Name = "testing"
	if err := rm.Update(r); err != nil {
		t.Fatal(err)
	}

	rid, err := rm.Get(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.ID, rid.ID)

	rname, err := rm.GetByName(rid.Name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.Name, rname.Name)

	all, err := rm.All(10, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, all, 1)

	if err := os.Remove(rm.RolePath()); err != nil {
		t.Fatal(err)
	}
}
