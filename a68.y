%{
// 28 june 2012
package main
%}
%start assembly

%term IDENT NUMBER CHARACTER STRING
%term ENCODING_NAME FUNCTION_NAME VARIABLE LABEL UNDEFINED_LABEL EQUATE
%term OPCODE EX_WORD EX_LONG DC_B DC_W DC_L DCB DCTO EQU
%term DATAREG ADDRESSREG SR CCR USP PC
%term INCLUDE INCBIN IF ELSE WHILE REPT FUNC RETURN MACRO PRINT ERROR CHARSET
%term TERM

%left OR
%left AND
%left EQ NE '<' LE '>' GE
%left '+' '-' '|' '^'
%left '*' '/' MOD '&' LSH RSH
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
		'(' ADDRESSREG ')'
	|	'(' ADDRESSREG ',' DATAREG EX_WORD ')'
	|	'(' ADDRESSREG ',' DATAREG EX_LONG ')'
	|	'(' PC ')'
	|	'(' PC ',' DATAREG EX_WORD ')'
	|	'(' PC ',' DATAREG EX_LONG ')'
	;

expression:
		IDENT
	|	VARIABLE
	|	LABEL
	|	UNDEFINED_LABEL
	|	EQUATE
	|	'@' IDENT
	|	NUMBER
	|	CHARACTER
	|	'(' expression ')'
	|	FUNCTION_NAME '(' ')'
	|	FUNCTION_NAME '(' exprlist ')'
	|	expression '+' expression
	|	expression '-' expression
	|	expression '*' expression
	|	expression '/' expression
	|	expression MOD expression
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
	|	DC_W exprlist
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
	|	ENCODING_NAME STRING
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
	;
