package object

import (
	"fmt"
	"testing"
)

func TestEnv(t *testing.T) {
	parent := NewEnvironment()
	child := CreateChildEnvironment(parent)
	fmt.Println(child.GetParentEnv())
}
