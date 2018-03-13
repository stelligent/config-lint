package filter

type Filter struct {
	Type      string
	Key       string
	Op        string
	Value     string
	Or        []Filter
	And       []Filter
	Not       []Filter
	ValueFrom FilterValueFrom
}

type FilterValueFrom struct {
	Bucket string
	Key    string
}

type Rule struct {
	Id       string
	Message  string
	Severity string
	Resource string
	Filters  []Filter
	Except   []string
	Tags     []string
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
