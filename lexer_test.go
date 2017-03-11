package SQLParser

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Lexer(t *testing.T){
	sqlStmt := "--this is a comment\nDROP TABLE IF EXISTS `user`;\n/* this is a comment */\nCREATE TABLE `customers` (\n  `id` bigint(60) NOT NULL AUTO_INCREMENT,\n  `username` varchar(20) DEFAULT NULL\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	//sqlStmt := "SELECT * FROM user"
	fmt.Printf("%q\n", sqlStmt)

	scan := NewScanner(strings.NewReader(sqlStmt))

	listOfTokens := []Tokens{
		ANNOTATION,
		DROP, WHITESPACE, TABLE, WHITESPACE, IF, WHITESPACE, EXISTS, WHITESPACE, IDENT, SEMI_COLON, WHITESPACE,
		ANNOTATION, WHITESPACE, CREATE, WHITESPACE, TABLE, WHITESPACE, IDENT, WHITESPACE, OPEN_PARENTH, WHITESPACE,
		IDENT, WHITESPACE, BIGINT, OPEN_PARENTH, SIZE, CLOSE_PARENTH, WHITESPACE, NOT, WHITESPACE, NULL, WHITESPACE, 
		AUTO_INCREMENT, COMMA, WHITESPACE,IDENT, WHITESPACE, VARCHAR, OPEN_PARENTH, SIZE, CLOSE_PARENTH, WHITESPACE, 
		DEFAULT, WHITESPACE, NULL, WHITESPACE, CLOSE_PARENTH, WHITESPACE, IDENT, EQUAL, IDENT, WHITESPACE, DEFAULT, 
		WHITESPACE, IDENT, EQUAL, IDENT, SEMI_COLON,
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

	for i := 0; i < len(tokens); i++ {
		if tokens[i] != listOfTokens[i] {
			t.Errorf("expected: %v found: %v", listOfTokens[i], tokens[i])
		}
	}

} 
