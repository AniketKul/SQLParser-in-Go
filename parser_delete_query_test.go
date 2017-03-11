package SQLParser_test

import (
	"reflect"
	"strings"
	"testing"
	"SQLParser"
)

func Test_DELETE_QueryParser(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *SQLParser.DeleteStatement
		err  string
	}{
		// Single field statement
		{
			s: `DELETE name FROM tbl`,
			stmt: &SQLParser.DeleteStatement{
				Fields:    []string{"name"},
				TableName: "tbl",
			},
		},

		// Multi-field statement
		{
			s: `DELETE first_name, last_name, age FROM my_table`,
			stmt: &SQLParser.DeleteStatement{
				Fields:    []string{"first_name", "last_name", "age"},
				TableName: "my_table",
			},
		},

		// delete all statement
		{
			s: `DELETE * FROM my_table`,
			stmt: &SQLParser.DeleteStatement{
				Fields:    []string{"*"},
				TableName: "my_table",
			},
		},

		// 
		{s: `foo`, err: `found "foo", expected DELETE`},
		{s: `DELETE !`, err: `found "!", expected field`},
		{s: `DELETE field xxx`, err: `found "xxx", expected FROM`},
		{s: `DELETE field FROM *`, err: `found "*", expected table name`},
		
	}

	for i, tt := range tests {
		stmt, err := SQLParser.NewParser(strings.NewReader(tt.s)).ParseDeleteStatements()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
	
}

