package dbpg

import (
	"fmt"
	"testing"
)

func TestPrepareUpdateRow(t *testing.T) {
	//Case 1
	queryStr := PrepareUpdateRow("public.sender", []string{"sender_type", "national_id_number", "tax_id_number"}, "")
	if queryStr != `UPDATE public.sender SET sender_type=$1, national_id_number=$2, tax_id_number=$3` {
		t.Errorf("NO: %v", queryStr)
		return
	}
	fmt.Printf("OK: %v\n", queryStr)
	//Case 2
	queryStr = PrepareUpdateRow("public.sender", []string{"sender_type", "national_id_number", "tax_id_number"}, "phone_number=$4")
	if queryStr != `UPDATE public.sender SET sender_type=$1, national_id_number=$2, tax_id_number=$3 WHERE phone_number=$4` {
		t.Errorf("NO: %v", queryStr)
	}
	fmt.Printf("OK: %v\n", queryStr)
}
