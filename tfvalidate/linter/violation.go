package linter

type Violation struct {
	Reason      string
	ResourceKey string
	Attribute   string
	Value       string
}
