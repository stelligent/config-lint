package linter

import (
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"strings"
)

func makeVarMap(variables []Variable) map[string]ast.Variable {
	m := map[string]ast.Variable{}
	for _, v := range variables {
		m["var."+v.Name] = ast.Variable{
			Type:  ast.TypeString,
			Value: v.Value.(string), // FIXME add error checking
		}
	}
	return m
}

func interpolate(s string, variables []Variable) string {
	fileFunc := ast.Function{
		ArgTypes:   []ast.Type{ast.TypeString},
		ReturnType: ast.TypeString,
		Variadic:   false,
		Callback: func(inputs []interface{}) (interface{}, error) {
			b, err := ioutil.ReadFile(inputs[0].(string))
			if err != nil {
				return "", nil
			}
			return strings.TrimSpace(string(b)), nil
		},
	}
	config := &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap: makeVarMap(variables),
			FuncMap: map[string]ast.Function{
				"file": fileFunc,
			},
		},
	}
	tree, err := hil.Parse(s)
	if err != nil {
		assertion.Debugf("Parse error: %v\n", err)
		return ""
	}
	result, err := hil.Eval(tree, config)
	if err != nil {
		assertion.Debugf("Eval error: %v\n", err)
		return ""
	}
	return result.Value.(string)
}
