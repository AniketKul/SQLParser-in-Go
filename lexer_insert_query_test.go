package SQLParser

import (
	"fmt"
	"strings"
	"testing"
)

func Test_INSERT_QueryLexer(t *testing.T){

	sqlStmt_q := "INSERT INTO Customers (CustomerName,ContactName,Address,City,PostalCode,Country) VALUES ('Cardinal','Tom B. Erichsen','Skagen 21','Stavanger','4006','Norway');"
	
	fmt.Printf("%q\n", sqlStmt_q)

	scan := NewScanner(strings.NewReader(sqlStmt_q))

	listOfTokens := []Tokens{
		INSERT, WHITESPACE, INTO, WHITESPACE, IDENT, WHITESPACE, OPEN_PARENTH, IDENT, COMMA, IDENT, COMMA, IDENT, COMMA, IDENT, COMMA, IDENT, COMMA, IDENT, CLOSE_PARENTH, WHITESPACE, VALUES, WHITESPACE, OPEN_PARENTH, STRING, COMMA, STRING, COMMA, STRING, COMMA, STRING, COMMA, STRING, COMMA, STRING, CLOSE_PARENTH, SEMI_COLON,
	}

	var tokens []Tokens

	for{
		if tok, litr :=scan.Scan(); tok!=EOF{
			fmt.Printf("%v: %v\n", tok, litr)
			tokens=append(tokens, tok)
		}else{
			break
		}
	}

	if len(tokens)!=len(listOfTokens){
		t.Errorf("Tokens Mismatch! expected %d but found %d\n", len(listOfTokens), len(tokens))
	}

	//fmt.Println(tokens)
	//fmt.Println(listOfTokens)

	for i := 0; i < len(tokens); i++ {
		if tokens[i] != listOfTokens[i] {
			t.Errorf("expected: %v found: %v", listOfTokens[i], tokens[i])
		}
	}

} 
