package jira

import (
	"fmt"
	"html"
	"strings"
)

func FormatMeta(meta map[string]interface{}) string {
	output := ""

	for k, v := range meta {
		value := ""
		switch v := v.(type) {
		default:
			break
		case bool:
			boolText := ""
			if v {
				boolText = "YES"
			} else {
				boolText = "NO"
			}
			value = fmt.Sprintf("|*%s*|%s|", k, boolText)
		case []interface{}:
			values := make([]string, len(v))
			for i, el := range v {
				values[i] = html.UnescapeString(fmt.Sprintf("%s", el))
			}
			value = fmt.Sprintf("|*%s*|%s|", k, strings.Join(values, "-"))
		case float64:
			value = fmt.Sprintf("|*%s*|%f|", k, v)
		case string:
			value = fmt.Sprintf("|*%s*|%s|", k, v)
		case map[string]interface{}:
			value = FormatMeta(v)
		}
		if value != "" {
			output += value + "\n"
		}
	}
	return output
}
