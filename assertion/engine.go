package assertion

type Assertion struct {
	Type      string
	Key       string
	Op        string
	Value     string
	Or        []Assertion
	And       []Assertion
	Not       []Assertion
	ValueFrom AssertionValueFrom `json:"value_from"`
}

type AssertionValueFrom struct {
	Url string
}

type ValueSource interface {
	GetValue(Assertion) string
}

type Rule struct {
	Id         string
	Message    string
	Severity   string
	Resource   string
	Assertions []Assertion
	Except     []string
	Tags       []string
}

type RuleSet struct {
	Type        string
	Description string
	Files       []string
	Rules       []Rule
	Version     string
}

type Violation struct {
	RuleId       string
	ResourceId   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
}

type ValidationReport struct {
	Violations   map[string]([]Violation)
	FilesScanned []string
}

type Resource struct {
	Id         string
	Type       string
	Properties interface{}
	Filename   string
}
