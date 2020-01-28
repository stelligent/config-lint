package linter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterFiles(t *testing.T) {
	input := []string {"bad.txt", "good.tf"}
	filters := [] string {"*.tf"}

	actual := filterFiles(input, filters)

	assert.Equal(t, 1, len(actual))
}
