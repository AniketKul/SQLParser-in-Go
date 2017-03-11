package SQLParser

import (
	"fmt"
	"strings"
	"testing"
)

func Test_DELETE_QueryLexer(t *testing.T){

	sqlStmt_del := "DELETE * FROM user"
	
	fmt.Printf("%q\n", sqlStmt_del)

	scan := NewScanner(strings.NewReader(sqlStmt_del))

	listOfTokens := []Tokens{
		DELETE, WHITESPACE, ASTERISK, WHITESPACE, FROM, WHITESPACE, IDENT, 
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
