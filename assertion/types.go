package assertion

type (

	// Resource describes a resource to be linted
	Resource struct {
		ID         string
		Type       string
		Properties interface{}
		Filename   string
		LineNumber int
	}

	// RuleSet describes a collection of rules for a Linter
	RuleSet struct {
		Type        string
		Description string
		Files       []string
		Rules       []Rule
		Version     string
		Resources   []ResourceConfig
	}

	// Rule is part of a RuleSet
	Rule struct {
		ID         string
		Message    string
		Severity   string
		Resource   string
		Conditions []Expression
		Assertions []Expression
		Except     []string
		Tags       []string
		Invoke     InvokeRuleAPI
	}

	// Expression expression for a Rule
	Expression struct {
		Key       string
		Op        string
		Value     string
		ValueType string    `json:"value_type"`
		ValueFrom ValueFrom `json:"value_from"`
		Or        []Expression
		Xor       []Expression
		And       []Expression
		Not       []Expression
		Every     CollectionExpression
		Some      CollectionExpression
		None      CollectionExpression
	}

	// CollectionExpression assertion for every element of a collection
	CollectionExpression struct {
		Key         string
		Expressions []Expression
	}

	// ValueFrom describes a external source for values
	ValueFrom struct {
		URL string
	}

	// InvokeRuleAPI describes an external API for linting a resource
	InvokeRuleAPI struct {
		URL     string
		Payload string
	}

	// ResourceConfig describes how to discover resouces in a YAML file
	ResourceConfig struct {
		ID   string
		Type string
		Key  string
	}

	// ValidationReport summarizes validation for resources using rules
	ValidationReport struct {
		FilesScanned     []string
		Violations       []Violation
		ResourcesScanned []ScannedResource
	}

	// Violation has details for a failed assertion
	Violation struct {
		RuleID           string
		ResourceID       string
		ResourceType     string
		Status           string
		RuleMessage      string
		AssertionMessage string
		Filename         string
		LineNumber       int
	}

	// ScannedResource has details for each resource scanned
	ScannedResource struct {
		ResourceID   string
		ResourceType string
		RuleID       string
		Status       string
		Filename     string
		LineNumber   int
	}

	// ValueSource interface to fetch dynamic values
	ValueSource interface {
		GetValue(Expression) (string, error)
	}

	// ExternalRuleInvoker defines an interface for invoking an external API
	ExternalRuleInvoker interface {
		Invoke(Rule, Resource) (string, []Violation, error)
	}

	// MatchResult has a true/false result, but also includes a message for better reporting
	MatchResult struct {
		Match   bool
		Message string
	}

	// Result returns a status, along with a message
	Result struct {
		Status  string
		Message string
	}
)
