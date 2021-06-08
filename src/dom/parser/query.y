%{
package parser
import "github.com/query-builder-generator/src/dom"
import "fmt"

type Token struct {
	token   int
	literal string
}
%}

%union{
    token Token
    query dom.Query
    identifier string
    classname string
}

%token <query> QUERY
%token <identifier> IDENTIFIER
%token FOR

%type <query> query
%type <classname> classname

%%

query      : QUERY IDENTIFIER FOR classname '{' '}'
        {
            $$ = dom.Query{Name: $2, Collection: $4}
            Domlex.(*Lexer).result = $$
        } ;

classname : IDENTIFIER
            {
                $$ = $1
            }
            | classname '.' IDENTIFIER
            {
                 $$ = fmt.Sprintf("%s.%s", $1, $3)
            } ;

%%