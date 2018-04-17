package parse

import (
	"fmt"
	"testing"
)

func TestComments(t *testing.T) {
	sql := `

SELECT CONCAT(aaaa,aaaa), SUM(coun) FROM table1 IGNORE INDEX (col3_index)
  WHERE col1=1 AND col2=2 AND col3=3;
`
	b := NewBuilder(sql)
	r, err := b.Parse()
	fmt.Println(r)
	fmt.Println(err)
}
