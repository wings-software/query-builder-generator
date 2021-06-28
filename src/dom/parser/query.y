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
    document dom.Document
    queries []dom.Query
    query dom.Query
    identifier string
    classname string
    filters []dom.Filter
    filter dom.Filter
    projectfields []string
    projectfield string
    token Token
}

%token <query> QUERY
%token <identifier> IDENTIFIER
%token EQUAL
%token FOR
%token IN
%token IS
%token LESS
%token LIST
%token MODULE
%token OF
%token PROJECT
%token REMAINDER
%token THAN

%type <document> document
%type <queries> query_list
%type <query> query
%type <classname> classname
%type <filters> filter_list
%type <filter> filter
%type <projectfields> projectfield_list
%type <projectfield> projectfield

%%

document    :   query_list
            {
                $$ = dom.Document{ Queries: $1 }
                Domlex.(*Lexer).result = $$
            }
            ;

query_list  : 	query
            {
                $$ = []dom.Query{$1}
            }
            |   query_list query
            {
                $$ = append($1, $2)
            }
            ;

query       :   QUERY IDENTIFIER FOR classname '{' filter_list '}'
            {
                $$ = dom.Query{Name: $2, Collection: $4, Filters: $6 }
                $$.Init()
            }
            |   QUERY IDENTIFIER FOR classname '{' filter_list projectfield_list '}'
            {
                $$ = dom.Query{Name: $2, Collection: $4, Filters: $6, ProjectFields: $7}
                $$.Init()
            }
            ;

filter_list : 	filter
            {
                $$ = []dom.Filter{$1}
            }
            |   filter_list filter
            {
                $$ = append($1, $2)
            }
            ;

filter      :   IDENTIFIER EQUAL classname ';'
            {
                $$ = dom.Filter{FieldType: $3, FieldName: $1, Operation: dom.Eq}
            }
            |   IDENTIFIER LESS THAN classname ';'
            {
                $$ = dom.Filter{FieldType: $4, FieldName: $1, Operation: dom.Lt}
            }
            |   IDENTIFIER IN LIST OF classname ';'
            {
                $$ = dom.Filter{FieldType: $5, FieldName: $1, Operation: dom.In}
            }
            |   IDENTIFIER MODULE REMAINDER IS ';'
            {
                $$ = dom.Filter{FieldName: $1, Operation: dom.Mod}
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

projectfield_list : projectfield
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
            ;
%%