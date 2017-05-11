package jira

import "fmt"

func FormatMeta(meta map[string]interface{}) string {
	output := ""

	for k, v := range meta {
		value := ""
		switch v := v.(type) {
		default:
			break
		case string:
			value = fmt.Sprintf("|*%s*|%s|", k, v)
		case map[string]interface{}:
			value = FormatMeta(v)
		}
		output += value + "\n"
	}
	return output
}
