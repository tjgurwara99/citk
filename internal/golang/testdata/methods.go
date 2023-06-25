package main

type MyType struct {
	camelCaseField             bool
	PascalCaseField            bool
	snake_case_field           bool
	SCREAMING_SNAKE_CASE_FIELD bool
	SCREAMINGFIELD             bool
}

func (m *MyType) camelCaseMethod()             {}
func (m *MyType) PascalCaseMethod()            {}
func (m *MyType) snake_case_method()           {}
func (m *MyType) SCREAMING_SNAKE_CASE_METHOD() {}
func (m *MyType) SCREAMINGMETHOD()             {}
