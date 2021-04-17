package translator

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// UnifiedRequest structure
type UnifiedRequest struct {
	Source string `json:"source"`
	// Modifiers will be here
	Requirements *RequirementExpression `json:"requirements"`
	Fields       []Field                `json:"fields"`
}

// RequirementExpression structure
type RequirementExpression struct {
	Type        string      `json:"type"`
	OrAnd       OrAnd       `json:"or_and"`
	Requirement Requirement `json:"requirement"`
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

// RequirementOperand structure
type RequirementOperand struct {
	Type string      `json:"type"`
	Spec interface{} `json:"spec"`
}

func (f Field) TransferToString() string {
	if f.Child != nil {
		return f.Name + "." + f.Child.TransferToString()
	} else {
		return f.Name
	}
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
	}

	return err
}
