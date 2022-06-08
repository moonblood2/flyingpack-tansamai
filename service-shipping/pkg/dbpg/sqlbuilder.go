package dbpg

import (
	"fmt"
	"strings"
)

func PrepareUpdateRow(table string, columns []string) string {
	colLen := len(columns) //columns length.
	query := fmt.Sprintf("UPDATE %v", table)
	setValues := make([]string, colLen)
	for i := 0; i < colLen; i++ {
		setValues[i] = fmt.Sprintf("%v=$%v", columns[i], i+1)
	}
	query = fmt.Sprintf("%v SET %v", query, strings.Join(setValues, ", "))
	return query
}
