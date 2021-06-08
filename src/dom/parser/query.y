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
    filters []dom.Filter
    filter dom.Filter
    projectfields []string
}

%token <query> QUERY
%token <identifier> IDENTIFIER
%token FOR
%token PROJECT
%token FILTER
%token AS
%token FROM
%token LIST

%type <query> query
%type <classname> classname
%type <filters> filter_list
%type <filter> filter
%type <projectfields> projectfields

%%

query     :	QUERY IDENTIFIER FOR classname '{' filter_list projectfields '}'
        	{
		    $$ = dom.Query{Name: $2, Collection: $4, Filters: $6, ProjectFields: $7}
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

filter_list: 	filter
		{
			$$ = []dom.Filter{$1}
		}
		| filter_list filter
		{
			$$ = append($1, $2)
		};

filter :	FILTER IDENTIFIER AS IDENTIFIER ';'
		{
		    $$ = dom.Filter{FieldType: $2, FieldName: $4, Operation: dom.EQUAL}
                }
                | FILTER IDENTIFIER AS IDENTIFIER FROM LIST ';'
		{
		    $$ = dom.Filter{FieldType: $2, FieldName: $4, Operation: dom.IN}
                } ;

projectfields: PROJECT IDENTIFIER ';'
		{
			$$ = []string{$2}
		}
%%