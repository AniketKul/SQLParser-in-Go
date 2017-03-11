package SQLParser

import "bufio"
import "strings"
import "bytes"
import "io"

type Scanner struct{
	r*bufio.Reader
}

//Let's declare tokens
type Tokens int

const(
	//special tokens
	ILLEGAL Tokens=iota	//0
	EOF     			//1
	ANNOTATION			//2
	WHITESPACE      	//3
	STRING      		//4
						// and so on
	//Literals
	IDENT

	//Special Characters
	COMMA
	ASTERISK
	BACKTICK
	SEMI_COLON
	OPEN_PARENTH
	CLOSE_PARENTH

	//Standard data types
	SIZE 
	BIT
	TINYINT
	SMALLINT
	INT
	BIGINT
	FLOAT
	DOUBLE
	LONGTEXT
	MEDIUMTEXT
	VARCHAR
	DATE
	TIME
	DATETIME
	TIMESTAMP

	//SQL keywords
	DROP
	LOCK
	UNLOCK
	TABLES
	WRITE
	IF
	EXISTS
	EQUAL
	CREATE
	TABLE
	DEFAULT
	NOT
	NULL
	COMMENT
	KEY
	UNIQUE
	CONSTRAINT
	PRIMARY
	FOREIGN
	REFERENCES
	AUTO_INCREMENT
	CURRENT_TIMESTAMP
	
	//SQL Query Keywords
	SELECT
	FROM
	INSERT
	INTO
	VALUES
	DELETE
	UPDATE
	SET
	WHERE
)

var (
eof = rune(0)
)

func isDigit(ch rune) bool {
	return (ch >='0' && ch<='9')
}

func isString(ch rune) bool {
	return ch=='\''
}

func isLetter(ch rune) bool {
	return ((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'))
}

func isWhiteSpace(ch rune) bool {
	return (ch == ' ' || ch == '\t' || ch == '\n')
}

//Create new scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (scan *Scanner) unread(){
	_ = scan.r.UnreadRune()
}

func (scan *Scanner) read() rune{
	ch, _, err := scan.r.ReadRune()

	if err!=nil {
		return eof
	}
	return ch
}

//This function is used to capture the digits in the expression. 
func (scan *Scanner) captureDigit() (tok Tokens, litr string) {
	var buf bytes.Buffer
	buf.WriteRune(scan.read())

	for{ 
		if ch := scan.read(); ch==eof {
			break
		}else if !isDigit(ch) {
			scan.unread()
			break
		}else{
			buf.WriteRune(ch)
		}
	}

	return SIZE, buf.String()
}

//This function is used to capture whitespaces in the expression.
func (scan *Scanner) captureWhiteSpace() (tok Tokens, litr string) {
	var buf bytes.Buffer
	buf.WriteRune(scan.read())

	for{
		if ch:=scan.read(); ch==eof{
			break
		}else if !isWhiteSpace(ch){
			scan.unread()
			break
		}else{
			buf.WriteRune(ch)
		}
	}

	return WHITESPACE, buf.String()
}

//This function will scan comments
func (scan *Scanner) scanComments() (tok Tokens, litr string){

	for{
		if ch:= scan.read(); ch==eof {
			return ILLEGAL, ""
		}else if ch=='*'{
			if c:= scan.read(); c=='/'{
				break
			}
		}
	}

	return ANNOTATION, ""

}

//This function will scan a string. 
func (scan *Scanner) scanString() (tok Tokens, litr string) {
	var buf bytes.Buffer
	ch := scan.read()

	readStr := func(c rune) {
		for {
			if ch := scan.read(); ch==c {
				break
			} else {
				_, _ = buf.WriteRune(ch)
			}
		}
	}

	switch ch {
	case '`':
		tok = IDENT
		readStr('`')
	case '\'':
		tok = STRING
		readStr('\'')
	default:
		return ILLEGAL, string(ch)
	}
	return tok, buf.String()
}

func (scan *Scanner) scanSQLKeyWords() (tok Tokens, litr string){
	var buf bytes.Buffer

	buf.WriteRune(scan.read())

	for{
		if ch := scan.read(); ch==eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			scan.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToUpper(buf.String()){

		case "DROP":
			return DROP, buf.String()
		case "IF":
			return IF, buf.String()
		case "EXISTS":
			return EXISTS, buf.String()
		case "LOCK":
			return LOCK, buf.String()
		case "UNLOCK":
			return UNLOCK, buf.String()
		case "TABLES":
			return TABLES, buf.String()
		case "WRITE":
			return WRITE, buf.String()
		case "CREATE":
			return CREATE, buf.String()
		case "TABLE":
			return TABLE, buf.String()
		case "NOT":
			return NOT, buf.String()
		case "NULL":
			return NULL, buf.String()
		case "DEFAULT":
			return DEFAULT, buf.String()
		case "COMMENT":
			return COMMENT, buf.String()
		case "KEY":
			return KEY, buf.String()
		case "UNIQUE":
			return UNIQUE, buf.String()
		case "CONSTRAINT":
			return CONSTRAINT, buf.String()
		case "PRIMARY":
			return PRIMARY, buf.String()
		case "FOREIGN":
			return FOREIGN, buf.String()
		case "REFERENCES":
			return REFERENCES, buf.String()
		case "AUTO_INCREMENT":
			return AUTO_INCREMENT, buf.String()
		case "CURRENT_TIMESTAMP":
			return CURRENT_TIMESTAMP, buf.String()
		case "BIT":
			return BIT, buf.String()
		case "TINYINT":
			return TINYINT, buf.String()
		case "SMALLINT":
			return SMALLINT, buf.String()
		case "INT":
			return INT, buf.String()
		case "BIGINT":
			return BIGINT, buf.String()
		case "FLOAT":
			return FLOAT, buf.String()
		case "DOUBLE":
			return DOUBLE, buf.String()
		case "VARCHAR":
			return VARCHAR, buf.String()
		case "LONGTEXT":
			return LONGTEXT, buf.String()
		case "MEDIUMTEXT":
			return MEDIUMTEXT, buf.String()
		case "DATE":
			return DATE, buf.String()
		case "TIME":
			return TIME, buf.String()
		case "DATETIME":
			return DATETIME, buf.String()
		case "TIMESTAMP":
			return TIMESTAMP, buf.String()

		//SQL Queries cases 
		case "SELECT":
			return SELECT, buf.String()
		case "FROM":
			return FROM, buf.String()
		case "INTO":
			return INTO, buf.String()
		case "VALUES":
			return VALUES, buf.String()
		case "INSERT":
			return INSERT, buf.String()
		case "DELETE":
			return DELETE, buf.String()
		case "UPDATE":
			return UPDATE, buf.String()
		case "SET":
			return SET, buf.String()
		case "WHERE":
			return WHERE, buf.String()

		default:
		return IDENT, buf.String()
	}	
}

func (scan *Scanner) Scan() (tok Tokens, litr string) {
	ch := scan.read()

	if isWhiteSpace(ch){
		scan.unread()
		return scan.captureWhiteSpace()
	}else if isLetter(ch){
		scan.unread()
		return scan.scanSQLKeyWords()
	} else if isDigit(ch) {
		scan.unread()
		return scan.captureDigit()
	} else if ch == '\'' || ch == '`' {
		scan.unread()
		return scan.scanString()
	} else if ch == '/' {
		if c := scan.read(); c == '*' {
			scan.unread()
			scan.unread()
			return scan.scanComments()
		}
		scan.unread()
		return ILLEGAL, string(ch)
	}

	switch ch {

	case eof:
		return EOF, "EOF"

	case ',':
		return COMMA, ","

	case '*':
		return ASTERISK, "*"

	case '(':
		return OPEN_PARENTH, "("

	case ')':
		return CLOSE_PARENTH, ")"

	case ';':
		return SEMI_COLON, ";"

	case '=':
		return EQUAL, "="

	case '-':
		if c := scan.read(); c == '-' { 
			for {
				if c := scan.read(); c == '\n' {
					return ANNOTATION, ""
				}
			}
		}
		return ILLEGAL, string(ch)
		
	default:
		return ILLEGAL, string(ch)
	}
}

