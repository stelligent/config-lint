package linter

import (
	"fmt"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"regexp"
	"strconv"
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
		m[v.Name] = makeVar(v.Value)
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

// The following functions are copied or adapted
// from https://github.com/hashicorp/terraform/blob/master/config/interpolate_funcs.go

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

// interpolationFuncJoin implements the "join" function that allows
// multi-variable values to be joined by some character.
func interpolationFuncJoin() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{ast.TypeString},
		Variadic:     true,
		VariadicType: ast.TypeList,
		ReturnType:   ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			var list []string

			if len(args) < 2 {
				return nil, fmt.Errorf("not enough arguments to join()")
			}

			for _, arg := range args[1:] {
				for _, part := range arg.([]ast.Variable) {
					if part.Type != ast.TypeString {
						return nil, fmt.Errorf(
							"only works on flat lists, this list contains elements of %s",
							part.Type.Printable())
					}
					list = append(list, part.Value.(string))
				}
			}

			return strings.Join(list, args[0].(string)), nil
		},
	}
}

// interpolationFuncConcat implements the "concat" function that concatenates
// multiple lists.
func interpolationFuncConcat() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{ast.TypeList},
		ReturnType:   ast.TypeList,
		Variadic:     true,
		VariadicType: ast.TypeList,
		Callback: func(args []interface{}) (interface{}, error) {
			var outputList []ast.Variable

			for _, arg := range args {
				for _, v := range arg.([]ast.Variable) {
					switch v.Type {
					case ast.TypeString:
						outputList = append(outputList, v)
					case ast.TypeList:
						outputList = append(outputList, v)
					case ast.TypeMap:
						outputList = append(outputList, v)
					default:
						return nil, fmt.Errorf("concat() does not support lists of %s", v.Type.Printable())
					}
				}
			}

			// we don't support heterogeneous types, so make sure all types match the first
			if len(outputList) > 0 {
				firstType := outputList[0].Type
				for _, v := range outputList[1:] {
					if v.Type != firstType {
						return nil, fmt.Errorf("unexpected %s in list of %s", v.Type.Printable(), firstType.Printable())
					}
				}
			}

			return outputList, nil
		},
	}
}

// interpolationFuncFormat implements the "format" function that does
// string formatting.
func interpolationFuncFormat() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{ast.TypeString},
		Variadic:     true,
		VariadicType: ast.TypeAny,
		ReturnType:   ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			format := args[0].(string)
			return fmt.Sprintf(format, args[1:]...), nil
		},
	}
}

// interpolationFuncList creates a list from the parameters passed
// to it.
func interpolationFuncList() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{},
		ReturnType:   ast.TypeList,
		Variadic:     true,
		VariadicType: ast.TypeAny,
		Callback: func(args []interface{}) (interface{}, error) {
			var outputList []ast.Variable

			for i, val := range args {
				switch v := val.(type) {
				case string:
					outputList = append(outputList, ast.Variable{Type: ast.TypeString, Value: v})
				case []ast.Variable:
					outputList = append(outputList, ast.Variable{Type: ast.TypeList, Value: v})
				case map[string]ast.Variable:
					outputList = append(outputList, ast.Variable{Type: ast.TypeMap, Value: v})
				default:
					return nil, fmt.Errorf("unexpected type %T for argument %d in list", v, i)
				}
			}

			// we don't support heterogeneous types, so make sure all types match the first
			if len(outputList) > 0 {
				firstType := outputList[0].Type
				for i, v := range outputList[1:] {
					if v.Type != firstType {
						return nil, fmt.Errorf("unexpected type %s for argument %d in list", v.Type, i+1)
					}
				}
			}

			return outputList, nil
		},
	}
}

// interpolationFuncReplace implements the "replace" function that does
// string replacement.
func interpolationFuncReplace() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeString, ast.TypeString, ast.TypeString},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			s := args[0].(string)
			search := args[1].(string)
			replace := args[2].(string)

			// We search/replace using a regexp if the string is surrounded
			// in forward slashes.
			if len(search) > 1 && search[0] == '/' && search[len(search)-1] == '/' {
				re, err := regexp.Compile(search[1 : len(search)-1])
				if err != nil {
					return nil, err
				}

				return re.ReplaceAllString(s, replace), nil
			}

			return strings.Replace(s, search, replace, -1), nil
		},
	}
}

// interpolationFuncElement implements the "element" function that allows
// a specific index to be looked up in a multi-variable value. Note that this will
// wrap if the index is larger than the number of elements in the multi-variable value.
func interpolationFuncElement() ast.Function {
	return ast.Function{
		ArgTypes:   []ast.Type{ast.TypeList, ast.TypeString},
		ReturnType: ast.TypeString,
		Callback: func(args []interface{}) (interface{}, error) {
			list := args[0].([]ast.Variable)
			if len(list) == 0 {
				return nil, fmt.Errorf("element() may not be used with an empty list")
			}

			index, err := strconv.Atoi(args[1].(string))
			if err != nil || index < 0 {
				return "", fmt.Errorf(
					"invalid number for index, got %s", args[1])
			}

			resolvedIndex := index % len(list)

			v := list[resolvedIndex]
			if v.Type != ast.TypeString {
				return nil, fmt.Errorf(
					"element() may only be used with flat lists, this list contains elements of %s",
					v.Type.Printable())
			}
			return v.Value, nil
		},
	}
}

// interpolationFuncMap creates a map from the parameters passed
// to it.
func interpolationFuncMap() ast.Function {
	return ast.Function{
		ArgTypes:     []ast.Type{},
		ReturnType:   ast.TypeMap,
		Variadic:     true,
		VariadicType: ast.TypeAny,
		Callback: func(args []interface{}) (interface{}, error) {
			outputMap := make(map[string]ast.Variable)

			if len(args)%2 != 0 {
				return nil, fmt.Errorf("requires an even number of arguments, got %d", len(args))
			}

			var firstType *ast.Type
			for i := 0; i < len(args); i += 2 {
				key, ok := args[i].(string)
				if !ok {
					return nil, fmt.Errorf("argument %d represents a key, so it must be a string", i+1)
				}
				val := args[i+1]
				variable, err := hil.InterfaceToVariable(val)
				if err != nil {
					return nil, err
				}
				// Enforce map type homogeneity
				if firstType == nil {
					firstType = &variable.Type
				} else if variable.Type != *firstType {
					return nil, fmt.Errorf("all map values must have the same type, got %s then %s", firstType.Printable(), variable.Type.Printable())
				}
				// Check for duplicate keys
				if _, ok := outputMap[key]; ok {
					return nil, fmt.Errorf("argument %d is a duplicate key: %q", i+1, key)
				}
				outputMap[key] = variable
			}

			return outputMap, nil
		},
	}
}

func Funcs() map[string]ast.Function {
	return map[string]ast.Function{
		"concat":  interpolationFuncConcat(),
		"element": interpolationFuncElement(),
		"file":    interpolationFuncFile(),
		"format":  interpolationFuncFormat(),
		"join":    interpolationFuncJoin(),
		"list":    interpolationFuncList(),
		"lookup":  interpolationFuncLookup(),
		"map":     interpolationFuncMap(),
		"replace": interpolationFuncReplace(),
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
		return s
	}
	result, err := hil.Eval(tree, config)
	if err != nil {
		assertion.Debugf("Eval error: %v\n", err)
		return s
	}
	return result.Value
}
