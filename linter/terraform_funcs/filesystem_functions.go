package terraform_funcs

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


var ReadFileContents = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "path",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		path := args[0].AsString()
		baseDir := filepath.Dir(path)

		fileByteContents, err := readFileBytes(baseDir, filepath.Base(path))
		if err != nil{
			return cty.NilVal, err
		}
		fileStringContent := strings.Trim(string(fileByteContents), "\n")
		return cty.StringVal(fileStringContent), nil
	},
})

func readFileBytes(baseDir, path string) ([]byte, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand ~: %s", err)
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}

	// Ensure that the path is canonical for the host OS
	path = filepath.Clean(path)

	src, err := ioutil.ReadFile(path)
	if err != nil {
		// ReadFile does not return Terraform-user-friendly error
		// messages, so we'll provide our own.
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no file exists at %s", path)
		}
		return nil, fmt.Errorf("failed to read %s", path)
	}

	return src, nil
}
