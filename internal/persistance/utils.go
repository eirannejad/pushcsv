package persistance

import (
	"fmt"
	"strings"
)

func ToSql(values *[]string) string {
	// wrap values in '' first
	cleanedValues := make([]string, 0)
	for _, value := range *values {
		cleanedValues = append(cleanedValues, fmt.Sprintf("'%s'", value))
	}
	// create the (,,,) sql value list
	return fmt.Sprintf("(%s)", strings.Join(cleanedValues, ", "))
}

func ToMap(fields, values *[]string) map[string]string {
	return make(map[string]string)
}
