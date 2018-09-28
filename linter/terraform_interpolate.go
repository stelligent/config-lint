package linter

import (
	"fmt"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"strings"
)

func makeList(variables []interface{}) []ast.Variable {
	list := []ast.Variable{}
	for _, v := range variables {
		list = append(list, makeVar(v))
	}
	return list
}

func makeMap(m map[string]interface{}) map[string]ast.Variable {
	result := map[string]ast.Variable{}
	for k, v := range m {
		if stringValue, ok := v.(string); ok {
			result[k] = ast.Variable{Type: ast.TypeString, Value: stringValue}
		}
	}
	return result
}

func makeVar(v interface{}) ast.Variable {
	switch tv := v.(type) {
	case string:
		return ast.Variable{
			Type:  ast.TypeString,
			Value: tv,
		}
	case []interface{}:
		return ast.Variable{
			Type:  ast.TypeList,
			Value: makeList(tv),
		}
	case map[string]interface{}:
		m := ast.Variable{
			Type:  ast.TypeMap,
			Value: makeMap(tv),
		}
		return m
	default:
		return ast.Variable{
			Type:  ast.TypeString,
			Value: "",
		}
	}
}

func makeVarMap(variables []Variable) map[string]ast.Variable {
	m := map[string]ast.Variable{}
	for _, v := range variables {
		m["var."+v.Name] = makeVar(v.Value)
	}
	return m
}

func interpolationFuncFile() ast.Function {
	return ast.Function{
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
}

// adapted from https://github.com/hashicorp/terraform/blob/master/config/interpolate_funcs.go
// interpolationFuncLookup implements the "lookup" function that allows
// dynamic lookups of map types within a Terraform configuration.
func interpolationFuncLookup() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{ast.TypeMap, ast.TypeString},
		ReturnType:   ast.TypeString,
		Variadic:     true,
		VariadicType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			defaultValue := ""
			defaultValueSet := false
			if len(args) > 2 {
				defaultValue = args[2].(string)
				defaultValueSet = true
			}
			if len(args) > 3 {
				return "", fmt.Errorf("lookup() takes no more than three arguments")
			}
			index := args[1].(string)
			mapVar := args[0].(map[string]ast.Variable)

			v, ok := mapVar[index]
			if !ok {
				if defaultValueSet {
					return defaultValue, nil
				} else {
					return "", fmt.Errorf(
						"lookup failed to find '%s'",
						args[1].(string))
				}
			}
			if v.Type != ast.TypeString {
				return nil, fmt.Errorf(
					"lookup() may only be used with flat maps, this map contains elements of %s",
					v.Type.Printable())
			}

			return v.Value.(string), nil
		},
	}
}

func Funcs() map[string]ast.Function {
	return map[string]ast.Function{
		"file":   interpolationFuncFile(),
		"lookup": interpolationFuncLookup(),
	}
}

func interpolate(s string, variables []Variable) interface{} {
	config := &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap:  makeVarMap(variables),
			FuncMap: Funcs(),
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
	return result.Value
}
