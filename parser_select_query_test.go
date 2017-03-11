package SQLParser_test

import (
	"reflect"
	"strings"
	"testing"
	"SQLParser"
)

func Test_SELECT_QueryParser(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *SQLParser.SelectStatement
		err  string
	}{
		// Single field statement
		{
			s: `SELECT name FROM tbl`,
			stmt: &SQLParser.SelectStatement{
				Fields:    []string{"name"},
				TableName: "tbl",
			},
		},

		// Multi-field statement
		{
			s: `SELECT first_name, last_name, age FROM my_table`,
			stmt: &SQLParser.SelectStatement{
				Fields:    []string{"first_name", "last_name", "age"},
				TableName: "my_table",
			},
		},

		// Select all statement
		{
			s: `SELECT * FROM my_table`,
			stmt: &SQLParser.SelectStatement{
				Fields:    []string{"*"},
				TableName: "my_table",
			},
		},

		// Errors
		{s: `foo`, err: `found "foo", expected SELECT`},
		{s: `SELECT !`, err: `found "!", expected field`},
		{s: `SELECT field xxx`, err: `found "xxx", expected FROM`},
		{s: `SELECT field FROM *`, err: `found "*", expected table name`},
	}

	for i, tt := range tests {
		stmt, err := SQLParser.NewParser(strings.NewReader(tt.s)).ParseSelectStatements()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
	
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}