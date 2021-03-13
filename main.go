package main

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
	Requirements RequirementExpression `json:"requirements"`
	Fields []Field `json:"fields"`
}

// UnifiedRequestToSql method. Transforms a UnifiedRequest structure into a full SQL statement.
func (u UnifiedRequest) UnifiedRequestToSql() string{
	var tmpRequest string
	var tmpFields string

	tmpRequest = "SELECT "

	for _,i := range u.Fields {
		tmpFields = i.TransferToString()
		tmpRequest += tmpFields + ", "
		tmpFields = ""
	}

	tmpRequest = tmpRequest[:len(tmpRequest) - 2] + " FROM " + u.Source
	tmpRequest = tmpRequest + " WHERE " + u.Requirements.ToSqlRequirementExpression()

	return tmpRequest
}

// RequirementExpression structure
type RequirementExpression struct {
	Type string `json:"type"`
	OrAnd OrAnd `json:"or_and"`
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
func (r RequirementExpression) ToSqlRequirementExpression() string  {
	var temp string
	switch r.Type {
	case "or":
		temp = "("
		for _, i := range r.OrAnd {
			temp += i.ToSqlRequirementExpression()
			temp += " OR "
		}
		return temp[:len(temp) - 4] + ")"

	case "and":
		temp = "("
		for _, i := range r.OrAnd {
			temp += i.ToSqlRequirementExpression()
			temp += " AND "
		}
		return temp[:len(temp) - 5] + ")"

	case "requirement":
		temp = "(" + r.Requirement.Field.TransferToString()
		switch r.Requirement.Operator {
		case "in": // Can be done with errors, but I'm not sure!
			temp += " IN ("
			for _,i := range r.Requirement.Value.Spec.([]interface{}) {
				temp += "'" + fmt.Sprint(i) + "', "
			}
			return temp[:len(temp) - 2] + "))"

		case "not in": // Can be done with errors, but I'm not sure!
			temp += " NOT IN ("
			for _,i := range r.Requirement.Value.Spec.([]interface{}) {
				temp += "'" + fmt.Sprint(i) + "', "
			}
			return temp[:len(temp) - 2] + "))"

		case "ne":
			return temp + " <> '" + fmt.Sprint(r.Requirement.Value.Spec) + "')"

		case "eq":
			return temp + " = '" + fmt.Sprint(r.Requirement.Value.Spec) + "')"

		case "ge":
			return temp + " >= " + fmt.Sprint(r.Requirement.Value.Spec) + ")"

		case "le":
			return temp + " <= " + fmt.Sprint(r.Requirement.Value.Spec) + ")"

		case "gt":
			return temp + " > " + fmt.Sprint(r.Requirement.Value.Spec) + ")"

		case "lt":
			return temp + " < " + fmt.Sprint(r.Requirement.Value.Spec) + ")"
		}
	}
	return ""
}

// OrAnd expression structure (transformed into slice)
type OrAnd []RequirementExpression


// Requirement structure
type Requirement struct {
	Operator string `json:"operator"`
	Field Field `json:"field"`
	Value RequirementOperand `json:"value"`
}

// Field identifier structure
type Field struct {
	Child *Field `json:"child"`
	Name string `json:"name"`
}

func (f Field) TransferToString() string  {

	if f.Child != nil {
		return f.Name + "." + f.Child.TransferToString()
	} else {
		return f.Name
	}
}

// RequirementOperand structure
type RequirementOperand struct {
	Type string `json:"type"`
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
func RequirementOperandDeclarator (r RequirementOperand) { //Error checking not done;Case raw_value not done
	switch r.Type {
	case "raw_value":


	case "request":
		tmp := r.Spec
		r.Spec = tmp.(Hook)
	}
}

// CleanGlobalVariables function. Makes global strings empty.

func main() {
	var input string
	fmt.Println("Choose the test! " +
		"\n -If you would like to test test1.json, please, enter 1" +
		"\n -If you would like to test test2.json, please, enter 2" +
		"\n -If you have your custom [request].json, please, enter c")
	_,_ = fmt.Scan(&input)
	switch input {
	case "1": Test1()
	case "2": Test2()
	case "c": CustomTest()
	default: fmt.Println("You've entered the wrong statement!")
	}
}

// Test1 function. Translates test1.json into an SQL statement and writes the result into the standard output
func Test1(){
	ConnectionPath1 := "C:\\Users\\tnche\\OneDrive\\Документы\\GitHub\\UnifiedRequestTranslator\\test1.json"

	fmt.Println("Activated test1.json...Press enter to continue")
	r := bufio.NewReader(os.Stdin)
	_,_,_ = r.ReadLine()

	var u UnifiedRequest
	u = u.Receive(ConnectionPath1)
	fmt.Println("Request translated into SQL:")
	fmt.Println("")
	fmt.Println(u.UnifiedRequestToSql())
}

// Test2 function. Translates test2.json into an SQL statement and writes the result into the standard output
func Test2(){
	ConnectionPath1 := "C:\\Users\\tnche\\OneDrive\\Документы\\GitHub\\UnifiedRequestTranslator\\test2.json"

	fmt.Println("Activated test2.json...Press enter to continue")
	r := bufio.NewReader(os.Stdin)
	_,_,_ = r.ReadLine()

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
