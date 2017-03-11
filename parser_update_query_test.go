package SQLParser_test

import (
	"reflect"
	"strings"
	"testing"
	"SQLParser"
)

func Test_UPDATE_QueryParser(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *SQLParser.UpdateStatement
		err  string
	}{
		// Single field statement
		{
			s: `UPDATE Customers SET City='Hamburg' WHERE CustomerID=1`,
			//s: `INSERT INTO Customers CustomerName,ContactName,Address,City,PostalCode,Country VALUES 'Cardinal','Tom B. Erichsen','Skagen 21','Stavanger','4006','Norway';`,
			stmt: &SQLParser.UpdateStatement{
				Fields:    []string{"City","Hamburg","CustomerID"},
				TableName: "Customers",
			},
		},
		
	}

	for i, tt := range tests {
		stmt, err := SQLParser.NewParser(strings.NewReader(tt.s)).ParseUpdateStatements()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
	
}

