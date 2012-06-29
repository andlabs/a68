%{
// 28 june 2012
package main
%}
%start assembly

%term IDENT NUMBER CHARACTER STRING
%term OPCODE EX_WORD EX_LONG DC_B DC_W DC_L DCB DCTO EQU
%term DATAREG ADDRESSREG SR CCR USP PC
%term INCLUDE INCBIN IF ELSE WHILE REPT FUNC RETURN MACRO PRINT ERROR CHARSET
%term TERM

%left OR
%left AND
%left EQ NE '<' LE '>' GE
%left '+' '-' '|' '^'
%left '*' '/' '%' '&' LSH RSH
%left UMINUS '~' '!'

%%
assembly:
		assembler_line		/* empty input handled within here */
	|	label_make assembly
	;

assembler_line:
		/* empty */
	|	operation TERM
	|	assignment TERM
	|	directive TERM
	;

label_make:
		IDENT ':'
	|	'@' IDENT ':'
	|	'+'
	|	'-'
	;

operation:
		OPCODE
	|	OPCODE operand
	|	OPCODE operand ',' operand
	|	data_definition
	;

operand:
		expression
	|	expression EX_WORD		/* TODO word and long keywords? */
	|	expression EX_LONG
	|	'#' expression
	|	DATAREG
	|	ADDRESSREG
	|	'-' '(' ADDRESSREG ')'
	|	'(' ADDRESSREG ')' '+'
	|	indirect
	|	expression indirect
	|	SR
	|	CCR
	|	USP
	|	'+'						/* TODO addressable? */
	|	'-'
	;

indirect:
		'(' address_or_pc ')'
	|	'(' address_or_pc ',' DATAREG EX_WORD ')'
	|	'(' address_or_pc ',' DATAREG EX_LONG ')'
	;

address_or_pc:
		ADDRESSREG
	|	PC
	;

expression:
		IDENT
	|	'@' IDENT
	|	NUMBER
	|	CHARACTER
	|	'(' expression ')'
	|	IDENT '(' exprlist ')'
	|	expression '+' expression
	|	expression '-' expression
	|	expression '*' expression
	|	expression '/' expression
	|	expression '%' expression
	|	expression '&' expression
	|	expression '|' expression
	|	expression '^' expression
	|	expression LSH expression
	|	expression RSH expression		/* TODO arithmetic */
	|	'-' expression					%prec UMINUS
	|	'~' expression
	|	expression EQ expression
	|	expression NE expression
	|	expression '<' expression
	|	expression LE expression
	|	expression '>' expression
	|	expression GE expression
	|	expression AND expression
	|	expression OR expression
	|	'!' expression
	;

data_definition:
		DC_B bytedatalist
		DC_W exprlist
	|	DC_L exprlist
	|	DCB ex_size expression ',' expression
	|	DCTO ex_size expression ',' expression
	;

ex_size:
		DC_B
	|	DC_W
	|	DC_L
	;

bytedatalist:
		bytedata
	|	bytedatalist ',' bytedata
	;

bytedata:
		expression
	|	STRING
	|	IDENT STRING
	;

exprlist:
		expression
	|	exprlist ',' expression
	;

assignment:
		varassign
	|	IDENT EQU expression
	;

varassign:
		IDENT '=' expression
	;

directive:
		INCLUDE STRING
	|	INCBIN STRING
	|	IF expression '{' assembly_no_labels '}'
	|	IF expression '{' assembly_no_labels '}' ELSE '{' assembly_no_labels '}'
	|	WHILE expression '{' assembly_no_labels '}'
	|	REPT expression '{' assembly_no_labels '}'
/*	|	function
	|	RETURN
	|	RETURN expression
	|	macro
*/	|	print
	|	charset
	;

assembly_no_labels:	/* TODO really have this? it's to avoid complications */
		/* empty */
	|	assembly_no_labels assembler_line
	;

function:
		FUNC IDENT '(' ')' '{' assembly_no_labels '}'
	|	FUNC IDENT '(' identlist ')' '{' assembly_no_labels '}'
	;

macro:
		MACRO IDENT '(' ')' '{' assembly '}'
	|	MACRO IDENT '(' identlist ')' '{' assembly '}'
	;

identlist:
		IDENT
	|	identlist ',' IDENT
	;

print:
		PRINT STRING
	|	PRINT STRING ',' exprlist
	|	ERROR STRING
	|	ERROR STRING ',' exprlist
	;

charset:
		CHARSET IDENT '{' charset_contents '}'
	;

charset_contents:
		charset_definition
	|	charset_contents charset_definition
	;

charset_definition:
		CHARACTER
	|	CHARACTER '=' exprlist
	|	CHARACTER EQ CHARACTER
	;

/* TODO STRUCTURES */
