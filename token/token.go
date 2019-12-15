// 4 december 2019
package token

import (
	"sync"
)

// Token is the set of lexical tokens in the assembler.
type Token int

// Predefined tokens.
// All their strings are string verions of the constant names.
const (
	// Special tokens.
	// These are not literals, operators, or keywords.
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Literal tokens.
	IDENT		// Used for non-keyword identifiers.

	nPredefined
)

type tokenType int
const (
	special tokenType = iota
	literal
	operator
	keyword
)

type token struct {
	typ		tokenType
	str		string
	prec		int
}

var tokens = make([]token, nPredefined, 256)
var tokensLock sync.RWMutex

func init() {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	tokens[ILLEGAL].typ = special
	tokens[ILLEGAL].str = "ILLEGAL"
	tokens[EOF].typ = special
	tokens[EOF].str = "EOF"
	tokens[COMMENT].typ = special
	tokens[COMMENT].str = "COMMENT"
	tokens[IDENT].typ = literal
	tokens[IDENT].str = "IDENT"
}

// AddLiteral adds a literal token type. str should be a constant name for the token type.
func AddLiteral(str string) Token {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	n := Token(len(tokens))
	tokens = append(tokens, token{
		typ:		literal,
		str:		str,
	})
	return n
}

// AddOperator adds an operator token type. str should be the operator text. If the operator is a binary operator, prec should be its precedence, with 1 being the lowest precedence, then 2, 3, 4, and so on. If this operator is not a binary operator, prec should be 0.
func AddOperator(str string, prec int) Token {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	n := Token(len(tokens))
	tokens = append(tokens, token{
		typ:		operator,
		str:		str,
		prec:		prec,
	})
	return n
}

var keywords = make(map[string]Token)
var keywordsLock sync.RWMutex

// AddKeyword adds a keyword token type. str should be the keyword itself. It panics if the keyword was already defined.
func AddKeyword(str string) Token {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	keywordsLock.Lock()
	defer keywordsLock.Unlock()
	if keywords[str] != 0 {
		panic("keyword " + str + " already defined")
	}
	n := Token(len(tokens))
	keywords[str] = n
	tokens = append(tokens, token{
		typ:		keyword,
		str:		str,
	})
	return n
}

// Lookup returns the token type for the identifier or keyword stored in str. If str does not store a keyword, IDENT is returned.
func Lookup(str string) Token {
	keywordsLock.RLock()
	defer keywordsLock.RUnlock()
	t, ok := keywords[str]
	if !ok {
		return IDENT
	}
	return t
}

// IsLiteral returns whether t is a literal.
func (t Token) IsLiteral() bool {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	return tokens[t].typ == literal
}

// IsOperator returns whether t is an operator.
func (t Token) IsOperator() bool {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	return tokens[t].typ == operator
}

// IsKeyword returns whether t is a keyword.
func (t Token) IsKeyword() bool {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	return tokens[t].typ == keyword
}

// Precedence returns the binary-operator precedence for t, as passed to AddOperator. If t does not represent an operator, 0 is returned.
func (t Token) Precedence() int {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	return tokens[t].prec
}

// String returns the string for t. The returned value is as described by the various Add functions.
func (t Token) String() string {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	return tokens[t].str
}
