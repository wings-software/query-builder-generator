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
    filters []string
    qprojections string
}

%token <query> QUERY
%token <identifier> IDENTIFIER
%token FOR
%token PROJECT
%token FILTER
%token AS

%type <query> query
%type <classname> classname
%type <filters> filters
%type <qprojections> qprojections

%%

query     :	QUERY IDENTIFIER FOR classname '{' filters PROJECT '{' qprojections '}' '}'
        	{
		    $$ = dom.Query{Name: $2, Collection: $4, QueryStmtType: $6, ProjectFields: $9}
		    Domlex.(*Lexer).result = $$
		} ;

classname :	IDENTIFIER
            	{
                    $$ = $1
            	}
            	| classname '.' IDENTIFIER
            	{
            	    $$ = fmt.Sprintf("%s.%s", $1, $3)
            	} ;

filters :	FILTER IDENTIFIER AS IDENTIFIER
		{
                    $$ = fmt.Sprintf("{%s %s}", $2, $4)
                } ;

qprojections :  IDENTIFIER
		{
			$$ = $1
		} ;

%%