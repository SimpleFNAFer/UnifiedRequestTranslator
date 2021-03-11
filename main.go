package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
		i.TransferToString(tmpFields)
		tmpRequest += tmpFields[1:] + ", "
		tmpFields = ""
	}

	tmpRequest = tmpRequest[:len(tmpRequest) - 2] + " FROM " + u.Source
	u.Requirements.ToSqlRequirementExpression()
	tmpRequest = tmpRequest + " WHERE " + SqlRequirementExpression

	defer CleanGlobalVariables()
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

// Global variables SqlRequirementExpression and FieldString.
// Help with recursive ToSqlRequirementExpression method.
var SqlRequirementExpression string
var FieldString string

// ToSqlRequirementExpression method. Transforms a RequirementExpression into an SQL condition (after WHERE)
func (r RequirementExpression) ToSqlRequirementExpression()  {
	switch r.Type {
	case "or":
		SqlRequirementExpression += "("
		for _, i := range r.OrAnd {
			i.ToSqlRequirementExpression()
			SqlRequirementExpression += " OR "
		}
		SqlRequirementExpression = SqlRequirementExpression[:len(SqlRequirementExpression) - 4] + ")"

	case "and":
		SqlRequirementExpression += "("
		for _, i := range r.OrAnd {
			i.ToSqlRequirementExpression()
			SqlRequirementExpression += " AND "
		}
		SqlRequirementExpression = SqlRequirementExpression[:len(SqlRequirementExpression) - 5] + ")"

	case "requirement":
		r.Requirement.Field.TransferToString(FieldString)
		FieldString = FieldString[1:]
		SqlRequirementExpression += "(" + FieldString
		switch r.Requirement.Operator {		//Idea collapse! Not enough information about data type of the FieldString!
		case "in":
			SqlRequirementExpression += " IN ("
			for _,i := range r.Requirement.Value.Spec.([]string) {
				SqlRequirementExpression += "'" + i + "', "
			}
			SqlRequirementExpression = SqlRequirementExpression[:len(SqlRequirementExpression) - 2] + ")"

		case "not in":
			SqlRequirementExpression += " NOT IN ("
			for _,i := range r.Requirement.Value.Spec.([]string) {
				SqlRequirementExpression += "'" + i + "', "
			}
			SqlRequirementExpression = SqlRequirementExpression[:len(SqlRequirementExpression) - 2] + ")"

		case "ne":
			SqlRequirementExpression += " <> '" + r.Requirement.Value.Spec.(string) + "'"

		case "eq":
			SqlRequirementExpression += " = '" + r.Requirement.Value.Spec.(string) + "'"

		case "ge":
			SqlRequirementExpression += " >= '" + r.Requirement.Value.Spec.(string) + "'"

		case "le":
			SqlRequirementExpression += " <= '" + r.Requirement.Value.Spec.(string) + "'"

		case "gt":
			SqlRequirementExpression += " > '" + r.Requirement.Value.Spec.(string) + "'"

		case "lt":
			SqlRequirementExpression += " < '" + r.Requirement.Value.Spec.(string) + "'"
		}
		SqlRequirementExpression += ")"
		FieldString = ""
	}
}

// OrAnd expression structure (transformed into slice)
type OrAnd []RequirementExpression


// Requirement structure
type Requirement struct {
	Operator string `json:"operator"`
	Field Field `json:"FieldString"`
	Value RequirementOperand `json:"value"`
}

// Field identifier structure
type Field struct {
	Child *Field `json:"child"`
	Name string `json:"name"`
}

func (f Field) TransferToString(s string)  {

	if f.Child != nil {
		f.Child.TransferToString(s)
	}
	s = "." + f.Name + s
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
func (u UnifiedRequest) Receive(b []byte) UnifiedRequest {
	err := json.Unmarshal(b, &u)
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
func CleanGlobalVariables (){
	SqlRequirementExpression = ""
	FieldString = ""
}

func main() {

}
