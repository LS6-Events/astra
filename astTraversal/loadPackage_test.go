package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
	"testing"
)

func TestLoadPackage(t *testing.T) {
	cachedPackages = make(map[string]*packages.Package)

	existingPkg := "time"
	nonExistingPkg := "github.com/user/project"

	t.Run("should load a package", func(t *testing.T) {
		pkg, err := LoadPackage(existingPkg, ".")
		assert.NoError(t, err)
		assert.NotNil(t, pkg)

		assert.Equal(t, existingPkg, pkg.PkgPath)

		assert.Len(t, cachedPackages, 1)
	})

	t.Run("should load a package from the cache", func(t *testing.T) {
		pkg, err := LoadPackage(existingPkg, ".")
		assert.NoError(t, err)
		assert.NotNil(t, pkg)

		assert.Equal(t, existingPkg, pkg.PkgPath)

		assert.Len(t, cachedPackages, 1)
	})

	t.Run("should return an error if the package does not exist", func(t *testing.T) {
		pkg, err := LoadPackage(nonExistingPkg, ".")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})
}
