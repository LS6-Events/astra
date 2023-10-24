package openapi

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollisionSafeKey(t *testing.T) {
	key := collisionSafeKey(astTraversal.NoBindingTag, "example", "package")
	expected := "package.example"
	assert.Equal(t, expected, key)
}

func TestGetPackageName(t *testing.T) {
	pkg := "github.com/ls6-events/astra"
	expected := "astra"
	name := getPackageName(pkg)
	assert.Equal(t, expected, name)
}

func TestMakeComponentRef(t *testing.T) {
	collisionSafeNames = make(map[string]string)

	name := "example"
	pkg := "package"
	collisionSafeNames[collisionSafeKey(astTraversal.NoBindingTag, name, pkg)] = "something.else"
	expected := "#/components/schemas/something.else"
	ref := makeComponentRef(astTraversal.NoBindingTag, name, pkg)
	assert.Equal(t, expected, ref)
}

func TestMakeComponentRefName(t *testing.T) {
	collisionSafeNames = make(map[string]string)

	name := "example"
	pkg := "package"
	collisionSafeNames[collisionSafeKey(astTraversal.NoBindingTag, name, pkg)] = "pkg.example"
	expected := "pkg.example"
	refName := makeComponentRefName(astTraversal.NoBindingTag, name, pkg)
	assert.Equal(t, expected, refName)
}

func TestMakeCollisionSafeNamesFromComponents(t *testing.T) {
	collisionSafeNames = make(map[string]string)

	// Initialize a sample list of astra.Field components
	components := []astra.Field{
		{
			Name:    "Field1",
			Package: "github.com/example/package1",
		},
		{
			Name:    "Field2",
			Package: "github.com/example/package2",
		},
		{
			Name:    "Field3",
			Package: "github.com/another/package1",
		},
	}

	// Call the function to generate collisionSafeNames
	makeCollisionSafeNamesFromComponents(components)

	// Define expected collision-safe names
	expectedNames := map[string]string{
		collisionSafeKey(astTraversal.NoBindingTag, "Field1", "github.com/example/package1"): "example.package1.Field1",
		collisionSafeKey(astTraversal.NoBindingTag, "Field2", "github.com/example/package2"): "package2.Field2",
		collisionSafeKey(astTraversal.NoBindingTag, "Field3", "github.com/another/package1"): "another.package1.Field3",
	}

	// Compare the generated collisionSafeNames with the expected names
	for key, expectedValue := range expectedNames {
		actualValue, ok := collisionSafeNames[key]
		if !ok {
			t.Errorf("Expected key %s not found in collisionSafeNames", key)
		} else if actualValue != expectedValue {
			assert.Equal(t, expectedValue, actualValue)
		}
	}

	// Clean up by resetting the collisionSafeNames map
	collisionSafeNames = make(map[string]string)
}
