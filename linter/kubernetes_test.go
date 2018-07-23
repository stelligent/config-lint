package linter

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestKubernetesSpecialVariables(t *testing.T) {
	loader := KubernetesResourceLoader{}
	filename := "./testdata/resources/pod.yml"
	loaded, err := loader.Load(filename)
	assert.Nil(t, err, "Expecting Load to not return an error")
	for _, resource := range loaded.Resources {
		properties := resource.Properties.(map[string]interface{})
		assert.Equal(t, filename, properties["__file__"])
		assert.Equal(t, filepath.Dir(filename), properties["__dir__"])
	}
}
