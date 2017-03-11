package SQLParser

import (
	"fmt"
	"io"
	"strconv"
)

type Column struct{
	Name string
	Type string 
	Size int 
	Default interface{}
	Comment string
	Nullable bool
	AutoIncr bool
}

type Constraint struct{
	Index string 
	ForeignKey string 
	TableName string 
	ColumnName string 
}

type Table struct{
	Name string
	Columns map[string]*Column
	PrimaryKey string 
	UniqueKeys map[string]string
	Keys map[string]string
	Constraints map[string]*Constraint
	Extras map[string]string
}

//schema is used to store table details
type Schema map[string]*Table

//Parser
type Parser struct{
	sc *Scanner 
	buf struct{
		tok Tokens
		litr string
		n int
	}
}

//Type stores SQL datatype tokens and their literal representation
var Type map[Tokens]string

func init(){
	Type = make(map[Tokens]string)
	Type[BIT]="bit"
	Type[TINYINT] = "tinyint"
	Type[SMALLINT] = "smallint"
	Type[INT] = "int"
	Type[BIGINT] = "bigint"
	Type[FLOAT] = "float"
	Type[DOUBLE] = "double"
	Type[VARCHAR] = "varchar"
	Type[LONGTEXT] = "longtext"
	Type[MEDIUMTEXT] = "mediumtext"
	Type[DATE] = "date"
	Type[TIME] = "time"
	Type[DATETIME] = "datetime"
	Type[TIMESTAMP] = "timestamp"
}

// NewParser returns a new parser for given reader
func NewParser(r io.Reader) *Parser {
	return &Parser{sc: NewScanner(r)}
}

func (p *Parser) scan()(tok Tokens, litr string){

	if p.buf.n!=0{
		p.buf.n=0
		return p.buf.tok, p.buf.litr 
	}

	tok, litr = p.sc.Scan()
	p.buf.tok, p.buf.litr = tok, litr
	return 
}

func (p *Parser) unScan(){
	p.buf.n=1
}

func (p *Parser) scanIgnoreWhiteSpace() (tok Tokens, litr string) {
	tok, litr = p.scan()
	if tok==WHITESPACE || tok == ANNOTATION {
		tok, litr = p.scan()
	}
	return
}

func (p *Parser) scanIdent()(tok Tokens, litr string){
	tok, litr = p.scanIgnoreWhiteSpace()

	if tok!= IDENT{
		return ILLEGAL, litr
	}

	return tok, litr
}

func (p *Parser) scanType()(string, int, error){
	tok, litr := p.scanIgnoreWhiteSpace()

	if tok>=BIT && tok<=TIMESTAMP{
		tok1, litr1 := p.scanIgnoreWhiteSpace()

		if tok1!=OPEN_PARENTH{
			p.unScan()
			return Type[tok], 0, nil
		}

		tok2, litr2 := p.scanIgnoreWhiteSpace()
		tok3, litr3 := p.scanIgnoreWhiteSpace()

		if tok2 != SIZE || tok3 != CLOSE_PARENTH{
			return "", 0, fmt.Errorf("found %q, expected type(integer)", litr+litr1+litr2+litr3)
		}

		size, _ := strconv.Atoi(litr2)
		return Type[tok], size, nil
	}

	return "", 0, fmt.Errorf("found %q, expected type", litr)
}

func (p *Parser) scanDefault()(string, error){
	tok, litr := p.scanIgnoreWhiteSpace()

	if tok!=DEFAULT{
		return "", fmt.Errorf("found %q, expected DEFUALT", litr)
	}

	tok, litr = p.scanIgnoreWhiteSpace()

	switch tok{
		case NULL: 
			return "null", nil
		case CURRENT_TIMESTAMP:
			return "Current Timestamp", nil
		case STRING:
			return litr, nil 
	}

	return "", fmt.Errorf("found %q, expected NULL or value", litr)
}

func (p *Parser) scanColumn() (*Column, error){
	var column = &Column{}
	tok, litr := p.scanIdent()

	if tok!=IDENT{
		return nil, fmt.Errorf("found %q, expected ident", litr)
	}

	column.Name = litr
	t, s, err := p.scanType()
	if err != nil {
		return nil, err
	}
	column.Type = t
	column.Size = s

	for{
		tok, litr = p.scanIgnoreWhiteSpace()

		switch tok{

			case DEFAULT: 
				p.unScan()
				val, err := p.scanDefault()
				if err!=nil{
					return nil, err
				}
				column.Default=val
				column.Nullable=val=="null"

			case NULL: 
				column.Nullable=true; 

			case NOT:
				tok1, litr1 := p.scanIgnoreWhiteSpace()

				if tok1!=NULL{
					return nil, fmt.Errorf("found %q, expected NULL", litr1)
				}
				column.Nullable=false

			case COMMENT:
				if tok1, litr1 :=p.scanIgnoreWhiteSpace(); tok1==STRING{
					column.Comment=litr1
				}else{
					return nil, fmt.Errorf("found %q, expected 'comment'", litr1)
				}

			case AUTO_INCREMENT:
				column.AutoIncr = true

			case COMMA, ASTERISK, CLOSE_PARENTH:
				p.unScan()
				return column, nil

			case EOF:
				return nil, fmt.Errorf("Unexpected EOF")

			default:
				return nil, fmt.Errorf("found %q, expected column constraint", litr)
		}
	}
}

func (p *Parser) scanPrimarykey()(string, error){
	tok1, litr1 := p.scanIgnoreWhiteSpace()
	tok2, litr2 := p.scanIgnoreWhiteSpace()

	if tok1!=PRIMARY || tok2!=KEY{
		return "", fmt.Errorf("found %q, expected PRIMARY KEY", litr1+litr2)
	}

	tok, litr := p.scanIgnoreWhiteSpace()

	if tok==OPEN_PARENTH{
		p.unScan()
		tok, litr = p.scanParenthIdent()

		if tok!=IDENT{
			return "", fmt.Errorf("found %q, expected ident", litr)
		}

		return litr, nil
	}

	tok, litr = p.scanIdent()

	if tok!=IDENT{
		return "", fmt.Errorf("found %q, expected ident", litr)
	}

	return litr, nil
}

func (p *Parser) scanParenthIdent() (Tokens, string){
	tok, litr := p.scanIgnoreWhiteSpace()

	if tok!=OPEN_PARENTH{
		return ILLEGAL, litr
	}

	tok, litr = p.scanIgnoreWhiteSpace()

	if tok==IDENT{
		tok1, litr1 := p.scanIgnoreWhiteSpace()
		if tok1!=CLOSE_PARENTH{
			return ILLEGAL, litr+litr1
		}

		return tok, litr
	}

	return ILLEGAL, ""
}

func (p *Parser) scanKey() (string, string, error){
	var index, column string
	tok, litr := p.scanIgnoreWhiteSpace()

	if tok!=KEY {
		return "", "", fmt.Errorf("found %q, expected KEY", litr)
	}

	//parse index 
	tok, litr = p.scanIgnoreWhiteSpace()
	if tok==IDENT{
		index=litr
	}else{
		return "", "", fmt.Errorf("found %q, expected index", litr)
	}

	//parse column 
	tok, litr = p.scanIgnoreWhiteSpace()

	if tok==IDENT{
		column=litr
	}else if tok==OPEN_PARENTH{
		p.unScan()
		tok, litr = p.scanParenthIdent()

		if tok!=IDENT{
			return "", "", fmt.Errorf("found %q, expected", litr)
		}
		column=litr
	}else{
		return "", "", fmt.Errorf("found %q, expected ident", litr)
	}

	return index, column, nil 
}

func (p *Parser) scanConstraint()(*Constraint, error){
	var constraint = &Constraint{}

	tok, litr := p.scanIgnoreWhiteSpace()

	if tok!=CONSTRAINT{
		return nil, fmt.Errorf("found %q, expected CONSTRAINT", litr)
	}

	tok, litr = p.scanIdent()

	if tok!=IDENT{
		return nil, fmt.Errorf("found %q, expected ident", litr)
	}

	constraint.Index=litr
	tok1, litr1 := p.scanIgnoreWhiteSpace()
	tok2, litr2 := p.scanIgnoreWhiteSpace()

	if tok1 != FOREIGN || tok2 != KEY {
		return nil, fmt.Errorf("found %q, expected FOREIGN KEY", litr1+litr2)
	}

	tok, litr = p.scanParenthIdent()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", litr)
	}

	constraint.ForeignKey=litr

	tok, litr=p.scanIgnoreWhiteSpace()
	if tok != REFERENCES{
		return nil, fmt.Errorf("found %q, expected REFERENCES", litr)
	}

	tok, litr = p.scanIdent()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected `table_name`", litr)
	}

	constraint.TableName = litr

	tok, litr = p.scanParenthIdent()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected (`column_name`)", litr)
	}

	constraint.ColumnName = litr

	return constraint, nil
}

func (p *Parser) scanKV() (string, string, error) {
	tok, litr := p.scanIgnoreWhiteSpace()
	tok1, litr1 := p.scanIgnoreWhiteSpace()
	tok2, litr2 := p.scanIgnoreWhiteSpace()

	if (tok != IDENT && tok != AUTO_INCREMENT) || tok1 != EQUAL || (tok2 != IDENT && tok2 != STRING && tok2 != SIZE) {
		return "", "", fmt.Errorf("found %q, expected key=value", litr+litr1+litr2)
	}
	return litr, litr2, nil
}

func (p *Parser) scanExtra() (map[string]string, error){
	extras := make(map[string]string)

	for{
		if tok, _ := p.scanIgnoreWhiteSpace(); tok!=SEMI_COLON{

			if tok != DEFAULT{
				p.unScan()
			}

			k,v,err := p.scanKV()

			if err!=nil {
				return nil, err
			}

			extras[k]=v
		}else{
			p.unScan()
			break
		}
	}

	return extras, nil 
}

//Parse one table
func (p *Parser) parse() (*Table, error){
	table := &Table{
		Columns: 		make(map[string]*Column),
		UniqueKeys:  	make(map[string]string),
		Keys:        	make(map[string]string),
		Constraints: 	make(map[string]*Constraint),
		Extras:      	make(map[string]string),
	}

	for{
		if tok, litr := p.scanIgnoreWhiteSpace(); tok==DROP || tok==LOCK || tok==UNLOCK || tok==ANNOTATION{
			for{
				if tok, _ := p.scanIgnoreWhiteSpace(); tok==SEMI_COLON{
					break
				}else if tok == EOF{
					return nil, nil
				}
			}
		}else if tok==SEMI_COLON || tok==ANNOTATION{
			continue
		}else if tok==CREATE{
			break
		}else if tok==EOF{
			return nil, nil
		}else{
			return nil, fmt.Errorf("unexpected %v: %q", tok, litr)
		}
	}

	if tok, litr := p.scanIgnoreWhiteSpace(); tok != TABLE {
		return nil, fmt.Errorf("found CREATE %q, expected CREATE TABLE", litr)
	}

	if tok, litr := p.scanIdent(); tok==IDENT{
		table.Name=litr
	}else{
		return nil, fmt.Errorf("found CREATE TABLE %d %q, expected CREATE TABLE `ident`", tok, litr)
	}

	//scan columns 
	if tok, litr := p.scanIgnoreWhiteSpace(); tok != OPEN_PARENTH{
		return nil, fmt.Errorf("found %q, expected (", litr)
	}

	for{
		tok, litr := p.scanIgnoreWhiteSpace()

		switch tok{

			case IDENT: 
				p.unScan()
				col, err := p.scanColumn()

				if err!=nil {
					return nil, err
				}
				table.Columns[col.Name]=col

			case PRIMARY:
				p.unScan()
				key, err := p.scanPrimarykey()

				if err!=nil {
					return nil, err
				}
				table.PrimaryKey=key

			case UNIQUE:
				k, v, err := p.scanKey()
				if err!=nil {
					return nil, err
				}
				table.UniqueKeys[k]=v

			case KEY:
				p.unScan()
				index, col, err := p.scanKey()
				if err!=nil{
					return nil, err
				}
				table.Keys[index]=col

			case CONSTRAINT:
				p.unScan()
				cos, err := p.scanConstraint()
				if err != nil {
					return nil, err
				}
				table.Constraints[cos.ForeignKey]=cos

			case CLOSE_PARENTH:
				tok, litr = p.scanIgnoreWhiteSpace()
				if tok != SEMI_COLON {
					p.unScan()
					extras, err := p.scanExtra()
					if err != nil {
						return nil, err
					}
					table.Extras = extras
				}
				return table, nil

			case COMMA:
				continue

			case ASTERISK:
				continue

			case SEMI_COLON:
				return table, nil

			default:
				return nil, fmt.Errorf("found %q, expected ident or primary or unique or key or constraint", litr)
		}	

	}
}

// Parse returns parsed table schema and an error
func (p *Parser) Parse() (Schema, error) {
	schema := make(Schema)
	for {
		table, err := p.parse()
		
		if err != nil {
			return schema, err // return already parsed tables and error
		}
		if table == nil { // parse done
			break
		}
		schema[table.Name] = table
	}
	return schema, nil
}






/* SQL Query Parsing */

type SelectStatement struct {
	Fields    []string
	TableName string
}

type InsertStatement struct {
	Fields    []string
	TableName string
}

type DeleteStatement struct {
	Fields    []string
	TableName string
}

type UpdateStatement struct {
	Fields    []string
	TableName string
}

// This function parses SQL SELECT statements.
func (p *Parser) ParseSelectStatements() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != COMMA {
			p.unScan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	// Return the successfully parsed statement.
	return stmt, nil
}

// This function parses SQL INSERT statements. 
func (p *Parser) ParseInsertStatements() (*InsertStatement, error) {
	stmtins := &InsertStatement{}

	//First token should be a "INSERT" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != INSERT {
		return nil, fmt.Errorf("found %q, expected INSERT", lit)
	}

	//Next keyword should be INTO
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}

	//Next keyword should denote table name.
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmtins.TableName = lit

	//Check if the column names start with '('
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != OPEN_PARENTH {
		return nil, fmt.Errorf("found %q, expected OPEN_PARENTH", lit)
	}

	//Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.

		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmtins.Fields = append(stmtins.Fields, lit)


		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != COMMA {
			p.unScan()
			break
		}
	}

	//Check if the column names end with ')'
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != CLOSE_PARENTH {
		return nil, fmt.Errorf("found %q, expected CLOSE_PARENTH", lit)
	}

	//After loop, next token should be "VALUES" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != VALUES {
		return nil, fmt.Errorf("found %q, expected VALUES", lit)
	}

	//Check if the string values start with '('
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != OPEN_PARENTH {
		return nil, fmt.Errorf("found %q, expected OPEN_PARENTH", lit)
	}

	//Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.

		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != STRING {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmtins.Fields = append(stmtins.Fields, lit)


		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != COMMA {
			p.unScan()
			break
		}
	}

	//Check if the string values end with ')'
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != CLOSE_PARENTH {
		return nil, fmt.Errorf("found %q, expected CLOSE_PARENTH", lit)
	}

	// Return the successfully parsed statement.
	return stmtins, nil
}

// This function parses SQL DELETE statements.
func (p *Parser) ParseDeleteStatements() (*DeleteStatement, error) {
	stmtdel := &DeleteStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != DELETE {
		return nil, fmt.Errorf("found %q, expected DELETE", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmtdel.Fields = append(stmtdel.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != COMMA {
			p.unScan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	//we should read the table name.
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmtdel.TableName = lit

	// Return the successfully parsed statement.
	return stmtdel, nil
}

// This function parses SQL UPDATE statements.
func (p *Parser) ParseUpdateStatements() (*UpdateStatement, error) {
	stmtupdate := &UpdateStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != UPDATE {
		return nil, fmt.Errorf("found %q, expected UPDATE", lit)
	}

	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmtupdate.TableName = lit

	// Next we should see the "SET" keyword.
	if tok, lit := p.scanIgnoreWhiteSpace(); tok != SET {
		return nil, fmt.Errorf("found %q, expected SET", lit)
	}

	for{
		// Read a field.

		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmtupdate.Fields = append(stmtupdate.Fields, lit)

		/*
		tok, lit := p.scanIgnoreWhiteSpace()
		if tok != EQUAL {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		*/
		
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != EQUAL {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}

		tok, lit = p.scanIgnoreWhiteSpace()
		if tok != STRING {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmtupdate.Fields = append(stmtupdate.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhiteSpace(); tok != COMMA {
			p.unScan()
			break
		}

	}
	
	tok, lit = p.scanIgnoreWhiteSpace()
	if tok != WHERE {
		return nil, fmt.Errorf("found %q, expected field", lit)
	}

	tok, lit = p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected field", lit)
	}
	stmtupdate.Fields = append(stmtupdate.Fields, lit)
	
	if tok, _ := p.scanIgnoreWhiteSpace(); tok != EQUAL {
			return nil, fmt.Errorf("found %q, expected field", lit)
	}

	tok, lit = p.scanIgnoreWhiteSpace()
	if tok != SIZE {
		return nil, fmt.Errorf("found %q, expected field", lit)
	}

	// Return the successfully parsed statement.
	return stmtupdate, nil
}
