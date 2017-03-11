package SQLParser

import (
	"fmt"
	"strings"
	"testing"
)

func Test_UPDATE_QueryLexer(t *testing.T){

	sqlStmt_q := "UPDATE customers SET city='Hamburg' WHERE customerId=1"
	
	fmt.Printf("%q\n", sqlStmt_q)

	scan := NewScanner(strings.NewReader(sqlStmt_q))

	listOfTokens := []Tokens{
		UPDATE, WHITESPACE, IDENT, WHITESPACE, SET, WHITESPACE, IDENT, EQUAL, STRING, WHITESPACE, WHERE, WHITESPACE, IDENT, EQUAL, SIZE,
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
