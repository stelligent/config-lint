package assertion

type (

	// Resource describes a resource to be linted
	Resource struct {
		ID         string
		Type       string
		Properties interface{}
		Filename   string
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
		Assertions []Assertion
		Except     []string
		Tags       []string
		Invoke     InvokeRuleAPI
	}

	// Assertion expression for a Rule
	Assertion struct {
		Key       string
		Op        string
		Value     string
		ValueType string    `json:"value_type"`
		ValueFrom ValueFrom `json:"value_from"`
		Or        []Assertion
		And       []Assertion
		Not       []Assertion
		Every     CollectionAssertion
		Some      CollectionAssertion
		None      CollectionAssertion
	}

	// CollectionAssertion assertion for every element of a collection
	CollectionAssertion struct {
		Key        string
		Assertions []Assertion
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
	}

	// ScannedResource has details for each resource scanned
	ScannedResource struct {
		ResourceID   string
		ResourceType string
		Status       string
	}

	// ValueSource interface to fetch dynamic values
	ValueSource interface {
		GetValue(Assertion) (string, error)
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
