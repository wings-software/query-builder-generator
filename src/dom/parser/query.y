%{
package parser
import "github.com/query-builder-generator/src/dom"

type Token struct {
	token   int
	literal string
}
%}

%union{
    token Token
    query dom.Query
    identifier string
}

%token <query> QUERY
%token <identifier> IDENTIFIER
%type <query> q1

%%

q1      : QUERY IDENTIFIER '{' '}'
        {
            $$ = dom.Query{Name: $2}
            Domlex.(*Lexer).result = $$
        } ;
%%