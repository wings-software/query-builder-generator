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
    projectfield string
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
%type <projectfields> projectfield_list
%type <projectfield> projectfield

%%

query       :   QUERY IDENTIFIER FOR classname '{' filter_list '}'
            {
                $$ = dom.Query{Name: $2, Collection: $4, Filters: $6 }
                $$.Init()
                Domlex.(*Lexer).result = $$
            }
            |   QUERY IDENTIFIER FOR classname '{' filter_list projectfield_list '}'
            {
                $$ = dom.Query{Name: $2, Collection: $4, Filters: $6, ProjectFields: $7}
                $$.Init()
                Domlex.(*Lexer).result = $$
            }
            ;

filter_list : 	filter
            	{
                	$$ = []dom.Filter{$1}
                }
            	| filter_list filter
            	{
            		$$ = append($1, $2)
            	}
		;

filter      :   FILTER IDENTIFIER AS classname ';'
		    {
			$$ = dom.Filter{FieldType: $4, FieldName: $2, Operation: dom.Eq}
            	}
            	| FILTER IDENTIFIER AS classname FROM LIST ';'
            	{
		        $$ = dom.Filter{FieldType: $4, FieldName: $2, Operation: dom.In}
            	}
            	;

classname   :   IDENTIFIER
            	{
                	$$ = $1
            	}
            	| classname '.' IDENTIFIER
            	{
                	$$ = fmt.Sprintf("%s.%s", $1, $3)
            	}
            	;

projectfield_list : 	projectfield
            	{
                	$$ = []string{$1}
                }
            	| projectfield_list projectfield
            	{
            		$$ = append($1, $2)
            	}
		;

projectfield: PROJECT IDENTIFIER ';'
		{
			$$ = $2
		}
%%