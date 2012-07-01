// 1 july 2012
package main

import (
	"strings"
)

// probably need to think of a better name
var validSymbolTypes := map[int]string{
	ENCODING_NAME:		"encoding",
	FUNCTION_NAME:		"function",
	VARAIBLE:			"variable",
	LABEL:				"label",
	UNDEFINED_LABEL:		"undefined (forward) label",
	EQUATE:				"equate",
}

type Symbol struct {
	Name	string		// TODO keep this?
	Type		int
	Value	uint32
	// TODO functions
}

type SymbolTable struct {
	m map[string]*Symbol
}

var Symbols, Locals SymbolTable

func init() {
	Symbols = NewSymbolTable()
	addBuiltins()
	Locals = NewSymbolTable()
}

func NewSymbolTable() SymbolTable {
	var s SymbolTable

	s.m = make(map[string]*Symbol)
	return s
}

func (s *SymbolTable) Add(name string, stype int) *Symbol {
	name = strings.ToLower(name)
	if existing, ok := s.m[name]; ok {
		FATAL_BUG("symbol %s defined: %v", name, existing)
	}
	if _, ok := validSymbolTypes[stype]; !ok {
		FATAL_BUG("attempt to create a symbol with invalid type %d", stype)
	}
	newSym := &Symbol{
		Name:	name,
		Type:	stype,
	}
	s.m[name] = newSym
	return newSym
}

func (s SymbolTable) Get(name string) *Symbol {
	name = strings.ToLower(name)
	sym, ok := s.m[name]
	if !ok {
		return nil
	}
	if ok && sym == nil {		// sanity check
		FATAL_BUG("nil symbol %s in symbol table", name)
	}
	return sym
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s type %s number value $%X\n",
		s.Name, validSymbolTypes[s.Type], s.Value)
}
