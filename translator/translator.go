package translator

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

// DisjunctiveNormalizer interface. Undone yet!!!
type DisjunctiveNormalizer interface {
	Normalize()
	CheckType() error
}

// UnifiedRequest structure
type UnifiedRequest struct {
	Source string `json:"source"`
	// Modifiers will be here
	Requirements *RequirementExpression `json:"requirements"`
	Fields       []Field                `json:"fields"`
}

// UnifiedRequestToSql method. Transforms a UnifiedRequest structure into a full SQL statement.
func (u UnifiedRequest) UnifiedRequestToSql() string {
	var tmpRequest string
	var tmpFields string

	tmpRequest = "SELECT "

	if len(u.Fields) == 0 {
		tmpFields = "*"
	} else {
		for _, i := range u.Fields {
			tmpFields += i.TransferToString()
			if i != u.Fields[len(u.Fields)-1] {
				tmpFields += ", "
			}

		}
	}

	tmpRequest += tmpFields + " FROM " + u.Source + " WHERE " + u.Requirements.ToSqlRequirementExpression()

	return tmpRequest
}

// RequirementExpression structure
type RequirementExpression struct {
	Type        string      `json:"type"`
	OrAnd       OrAnd       `json:"or_and"`
	Requirement Requirement `json:"requirement"`
}

// CheckType method. Declares type for the "spec" FieldString
func (r RequirementExpression) CheckType() error {
	var err error

	switch r.Type {
	case "or":
		tmp := r.OrAnd
		if len(tmp) < 2 {
			err = errors.New("CheckType(): incorrect number of operands in or-requirement expression (expected > 1)")
		}
	case "and":
		tmp := r.OrAnd
		if len(tmp) < 2 {
			err = errors.New("CheckType(): incorrect number of operands in and-requirement expression (expected > 1)")
		}
	case "requirement":
		tmp := r.Requirement
		RequirementOperandDeclarator(tmp.Value)
	}

	return err
}

// Normalize method. Undone yet!!!
func (r RequirementExpression) Normalize() RequirementExpression {
	var n RequirementExpression

	if (r.Type == "or") || (r.Type == "and") {
		for _, i := range r.OrAnd {
			i.Normalize()
		}
	}
	fmt.Println(r.Type)

	return n
}

// ToSqlRequirementExpression method. Transforms a RequirementExpression into an SQL condition (after WHERE)
func (r RequirementExpression) ToSqlRequirementExpression() string {
	var temp string
	switch r.Type {
	case "or":
		temp = "(" + r.OrAnd[0].ToSqlRequirementExpression()
		for _, i := range r.OrAnd[1:] {
			temp += " OR " + i.ToSqlRequirementExpression()
		}
		return temp + ")"

	case "and":
		temp = "(" + r.OrAnd[0].ToSqlRequirementExpression()
		for _, i := range r.OrAnd[1:] {
			temp += " AND " + i.ToSqlRequirementExpression()
		}
		return temp + ")"

	case "requirement":
		return r.Requirement.ToSqlRequirement()
	}
	return ""
}

// ToSqlRequirement method. Translates a requirement structure into an SQL requirement
func (r Requirement) ToSqlRequirement() string {
	var temp string
	temp = "(" + r.Field.TransferToString()
	switch r.Operator {
	case "in": // Can be done with errors, but I'm not sure!
		values := r.Value.Spec.([]interface{})
		temp += " IN (" + "'" + fmt.Sprint(values[0]) + "'"
		for _, i := range values[1:] {
			temp += ", " + "'" + fmt.Sprint(i) + "'"
		}
		return temp + "))"

	case "not in": // Can be done with errors, but I'm not sure!
		values := r.Value.Spec.([]interface{})
		temp += " NOT IN (" + "'" + fmt.Sprint(values[0]) + "'"
		for _, i := range values[1:] {
			temp += ", " + "'" + fmt.Sprint(i) + "'"
		}
		return temp + "))"

	case "ne":
		return temp + " <> '" + fmt.Sprint(r.Value.Spec) + "')"

	case "eq":
		return temp + " = '" + fmt.Sprint(r.Value.Spec) + "')"

	case "ge":
		return temp + " >= " + fmt.Sprint(r.Value.Spec) + ")"

	case "le":
		return temp + " <= " + fmt.Sprint(r.Value.Spec) + ")"

	case "gt":
		return temp + " > " + fmt.Sprint(r.Value.Spec) + ")"

	case "lt":
		return temp + " < " + fmt.Sprint(r.Value.Spec) + ")"
	}

	return ""
}

// OrAnd expression structure (transformed into slice)
type OrAnd []RequirementExpression

// Requirement structure
type Requirement struct {
	Operator string             `json:"operator"`
	Field    Field              `json:"field"`
	Value    RequirementOperand `json:"value"`
}

// Field identifier structure
type Field struct {
	Child *Field `json:"child"`
	Name  string `json:"name"`
}

func (f Field) TransferToString() string {
	if f.Child != nil {
		return f.Name + "." + f.Child.TransferToString()
	} else {
		return f.Name
	}
}

// RequirementOperand structure
type RequirementOperand struct {
	Type string      `json:"type"`
	Spec interface{} `json:"spec"`
}

// Hook request structure
type Hook struct {
	// Undone yet!!!
}

// Receive method. Takes the whole unified request
func (u UnifiedRequest) Receive(path string) UnifiedRequest {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// RequirementOperandDeclarator method. Undone yet!!!
func RequirementOperandDeclarator(r RequirementOperand) { //Error checking not done;Case raw_value not done
	switch r.Type {
	case "raw_value":

	case "request":
		tmp := r.Spec
		r.Spec = tmp.(Hook)
	}
}

// CleanGlobalVariables function. Makes global strings empty.

// Test1 function. Translates test1.json into an SQL statement and writes the result into the standard output
func Test1() {
	ConnectionPath1 := "C:\\Users\\tnche\\OneDrive\\Документы\\GitHub\\UnifiedRequestTranslator\\test1.json"

	fmt.Println("Activated test1.json...Press enter to continue")
	r := bufio.NewReader(os.Stdin)
	_, _, _ = r.ReadLine()

	var u UnifiedRequest
	u = u.Receive(ConnectionPath1)
	fmt.Println("Request translated into ES:")
	fmt.Println("")
	b, _ := json.Marshal(u.UnifiedRequestToES())
	fmt.Println(string(b))
}

// Test2 function. Translates test2.json into an SQL statement and writes the result into the standard output
func Test2() {
	ConnectionPath1 := "C:\\Users\\tnche\\OneDrive\\Документы\\GitHub\\UnifiedRequestTranslator\\test2.json"

	fmt.Println("Activated test2.json...Press enter to continue")
	r := bufio.NewReader(os.Stdin)
	_, _, _ = r.ReadLine()

	var u UnifiedRequest
	u = u.Receive(ConnectionPath1)
	fmt.Println("Request translated into SQL:")
	fmt.Println("")
	fmt.Println(u.UnifiedRequestToSql())
}

// CustomTest function. Allows you to write a connection string manually
func CustomTest() {
	fmt.Println("Write down path to .json file:")
	var ConnectionPath string
	_, _ = fmt.Scan(&ConnectionPath)
	fmt.Println("Path received: " + ConnectionPath)
	var u UnifiedRequest
	u = u.Receive(ConnectionPath)
	fmt.Println("Request translated into SQL:")
	fmt.Println("")
	fmt.Println(u.UnifiedRequestToSql())
}

/*
es:
{
  "_source": [
    "name",
  ],
  "query": {
    "term": {
      "name": "noname"
    }
  }
}
sql:
select name where name="noname"
*/

/*
es:
{
  "_source": [
    "name",
  ],

  "query": {
    "bool": {
      "should": [
        {
          "term": {
            "name": "noname"
          }
        },
	{
	  "term": {
	    "surname": "petrov"
	  }
	}
      ]
    }
  }
}
sql:
select name where name="noname" or surname="petrov"
*/

type ESRequest struct {
	Source []string `json:"_source,omitempty"`
	Query  Query    `json:"query"`
}

type Query struct {
	Term     map[string]interface{}     `json:"term,omitempty"`
	Terms    map[string][]interface{}   `json:"terms,omitempty"`
	Bool     *Bool                      `json:"bool,omitempty"`
	Range    map[string]RangeExpression `json:"range,omitempty"`
	MatchAll map[string]interface{}     `json:"match_all,omitempty"`
}

type RangeExpression struct {
	Gte interface{} `json:"gte,omitempty"`
	Gt  interface{} `json:"gt,omitempty"`
	Lt  interface{} `json:"lt,omitempty"`
	Lte interface{} `json:"lte,omitempty"`
}

type Bool struct {
	// Из примеров можно увидеть что поля should содержат список Query
	Should []Query `json:"should,omitempty"`
	Must   []Query `json:"must,omitempty"`
}

func (u UnifiedRequest) UnifiedRequestToES() ESRequest {
	var tmpRequest ESRequest

	// Это неверно, если длина полей не равна нулю, то ими заполняется tmpRequest.Source
	// MatchAll тут ни при чем
	if u.Requirements == nil {
		tmpRequest.Query.MatchAll = make(map[string]interface{})
	} else {
		tmpRequest.Query = u.Requirements.ToESRequirementExpression()
	}

	if len(u.Fields) > 0 {
		tmpRequest.Source = make([]string, len(u.Fields))
		for i := 0; i < len(u.Fields); i++ {
			tmpRequest.Source[i] = u.Fields[i].TransferToString()
		}
	}

	return tmpRequest
}

func (r RequirementExpression) ToESRequirementExpression() Query {
	var temp Query

	switch r.Type {
	case "or":
		for _, i := range r.OrAnd {
			temp.Bool.Should = append(temp.Bool.Should, i.ToESRequirementExpression())
		}
		return temp

	case "and":
		for _, i := range r.OrAnd {
			temp.Bool.Must = append(temp.Bool.Must, i.ToESRequirementExpression())
		}
		return temp

	case "requirement":
		return r.Requirement.ToESRequirement()
	}
	return Query{}
}

//Doing now
func (r Requirement) ToESRequirement() Query {
	var temp Query

	switch r.Operator {
	case "in":
		temp.Terms = make(map[string][]interface{})

	/*case "not in":
		temp.Range = make(map[string]RangeExpression)

	case "ne":*/

	case "eq":
		temp.Term = map[string]interface{}{r.Field.TransferToString(): r.Value.Spec}

	case "ge":
		temp.Range = map[string]RangeExpression{r.Field.TransferToString(): RangeExpression{
			Gte: r.Value.Spec,
		}}

	case "le":
		temp.Range = map[string]RangeExpression{r.Field.TransferToString(): RangeExpression{
			Lte: r.Value.Spec,
		}}

	case "gt":
		temp.Range = map[string]RangeExpression{r.Field.TransferToString(): RangeExpression{
			Gt: r.Value.Spec,
		}}

	case "lt":
		temp.Range = map[string]RangeExpression{r.Field.TransferToString(): RangeExpression{
			Lt: r.Value.Spec,
		}}

	}

	return temp
}
