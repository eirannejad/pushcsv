package persistance

import (
	"fmt"
	"strings"
)

func ToSql(values *[]string) string {
	return fmt.Sprintf("( %s )", strings.Join(*values, ", "))
}
