package SQLParser_test

import (
	"reflect"
	"strings"
	"testing"
	"SQLParser"
)

func Test_INSERT_QueryParser(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *SQLParser.InsertStatement
		err  string
	}{
		
		{
			s: `INSERT INTO Customers (CustomerName,ContactName,Address,City,PostalCode,Country) VALUES ('Cardinal','Tom B. Erichsen','Skagen 21','Stavanger','4006','Norway');`,
			//s: `INSERT INTO Customers CustomerName,ContactName,Address,City,PostalCode,Country VALUES 'Cardinal','Tom B. Erichsen','Skagen 21','Stavanger','4006','Norway';`,
			stmt: &SQLParser.InsertStatement{
				Fields:    []string{"CustomerName","ContactName","Address","City","PostalCode","Country","Cardinal","Tom B. Erichsen","Skagen 21","Stavanger","4006","Norway"},
				TableName: "Customers",
			},
		},


	}

	for i, tt := range tests {
		stmt, err := SQLParser.NewParser(strings.NewReader(tt.s)).ParseInsertStatements()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
	
}

