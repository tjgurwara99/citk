package annotation

import (
	"fmt"
	"strings"
)

type AnnotationType string

const (
	Notice  AnnotationType = "notice"
	Warning AnnotationType = "warning"
	Error   AnnotationType = "error"
	Debug   AnnotationType = "debug"
)

type Annotation struct {
	FileName string
	Title    string
	Message  string
	// lines are of type uint32 because tree-sitter uses this for StartPoint().Row
	StartLine uint32
	EndLine   uint32
	StartCol  uint32
	EndCol    uint32
	Type      AnnotationType
}

func (a Annotation) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("::%s ", a.Type))
	if a.FileName != "" {
		builder.WriteString(fmt.Sprintf("file=%s,", a.FileName))
	}
	if a.StartLine != 0 {
		builder.WriteString(fmt.Sprintf("line=%d,", a.StartLine))
	}
	if a.EndLine != 0 {
		builder.WriteString(fmt.Sprintf("endLine=%d,", a.EndLine))
	}
	if a.StartCol != 0 {
		builder.WriteString(fmt.Sprintf("col=%d,", a.StartCol))
	}
	if a.EndCol != 0 {
		builder.WriteString(fmt.Sprintf("endCol=%d", a.EndCol))
	}
	if a.Message != "" {
		builder.WriteString("::" + a.Message)
	}
	return builder.String()
}
