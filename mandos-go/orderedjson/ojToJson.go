package orderedjson

import (
	"fmt"
	"strings"
)

// JSONString returns a formatted string representation of an ordered JSON
func JSONString(j OJsonObject) string {
	var sb strings.Builder
	j.writeJSON(&sb, 0)
	return sb.String()
}

func addIndent(sb *strings.Builder, indent int) {
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}
}

func (j *OJsonMap) writeJSON(sb *strings.Builder, indent int) {
	if j.Size() == 0 {
		sb.WriteString("{}")
		return
	}

	sb.WriteString("{")
	for i, child := range j.OrderedKV {
		sb.WriteString("\n")
		addIndent(sb, indent+1)
		sb.WriteString("\"")
		sb.WriteString(child.Key)
		sb.WriteString("\": ")
		if child.Value != nil {
			child.Value.writeJSON(sb, indent+1)
		}
		if i < len(j.OrderedKV)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("\n")
	addIndent(sb, indent)
	sb.WriteString("}")
}

func (j *OJsonList) writeJSON(sb *strings.Builder, indent int) {
	collection := j.AsList()
	if len(collection) == 0 {
		sb.WriteString("[]")
		return
	}

	sb.WriteString("[")
	for i, child := range collection {
		sb.WriteString("\n")
		addIndent(sb, indent+1)
		child.writeJSON(sb, indent+1)
		if i < len(collection)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("\n")
	addIndent(sb, indent)
	sb.WriteString("]")
}

func (j *OJsonString) writeJSON(sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("\"%s\"", j.Value))
}

func (j *OJsonBool) writeJSON(sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("%v", bool(*j)))
}
