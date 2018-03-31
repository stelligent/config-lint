package assertion

// Assertion expression for a resource
type Assertion struct {
	Key       string
	Op        string
	Value     string
	ValueType string `json:"value_type"`
	Or        []Assertion
	And       []Assertion
	Not       []Assertion
	ValueFrom ValueFrom `json:"value_from"`
}

// ValueFrom describes source for external values
type ValueFrom struct {
	URL string
}

// ValueSource interface to fetch values
type ValueSource interface {
	GetValue(Assertion) (string, error)
}

// InvokeRuleAPI describes parameters for calling an external API
type InvokeRuleAPI struct {
	URL     string
	Payload string
}

// Rule for a resource
type Rule struct {
	ID         string
	Message    string
	Severity   string
	Resource   string
	Assertions []Assertion
	Except     []string
	Tags       []string
	Invoke     InvokeRuleAPI
}

// RuleSet describes a collection of rules for a Linter
type RuleSet struct {
	Type        string
	Description string
	Files       []string
	Rules       []Rule
	Version     string
	Resources   []ResourceConfig
}

// ResourceConfig describes how to discover resouces in a YAML file
type ResourceConfig struct {
	ID   string
	Type string
	Key  string
}

// Violation has details for a failed assertion
type Violation struct {
	RuleID       string
	ResourceID   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
}

// ValidationReport summarizes validation for resources and rules
type ValidationReport struct {
	Violations   map[string]([]Violation)
	FilesScanned []string
}

// Resource describes a resource to be validated
type Resource struct {
	ID         string
	Type       string
	Properties interface{}
	Filename   string
}

// ExternalRuleInvoker defines an interface for invoking an external API
type ExternalRuleInvoker interface {
	Invoke(Rule, Resource) (string, []Violation, error)
}
