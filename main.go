package main

import (
	"encoding/json"
	"errors"
	"log"
)

// UnifiedRequest structure
type UnifiedRequest struct {
	Source string `json:"source"`
	// Modifiers will be here
	Requirements []RequirementExpression `json:"requirements"`
	Fields []Field `json:"fields"`
}

// RequirementExpression structure
type RequirementExpression struct {
	Type string `json:"type"`
	Spec interface{} `json:"spec"`
}

// OrAnd expression structure transformed into slice
type OrAnd []RequirementExpression


// Requirement structure
type Requirement struct {
	Operator string `json:"operator"`
	Field Field `json:"field"`
	Value RequirementOperand `json:"value"`
}

// Field identifier structure
type Field struct {
	Parent *Field `json:"parent"`
	Name string `json:"name"`
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

// UnifiedRequestHandler method. Takes the whole unified request
func UnifiedRequestHandler(b []byte) UnifiedRequest {
	var u UnifiedRequest
	err := json.Unmarshal(b, &u)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// RequirementExpressionDeclarator method. Declares type for the "spec" field
func RequirementExpressionDeclarator (r RequirementExpression) error {
	var err error

	switch r.Type {
	case "or":
		tmp := r.Spec
		if len(tmp.(OrAnd)) < 2 {
			err = errors.New("incorrect or-request")
		} else {
			r.Spec = tmp.(OrAnd)
		}
	case "and":
		tmp := r.Spec
		if len(tmp.(OrAnd)) < 2 {
			err = errors.New("incorrect and-request")
		} else {
			r.Spec = tmp.(OrAnd)
		}
	case "requirement":
		tmp := r.Spec
		RequirementOperandDeclarator(tmp.(Requirement).Value)
		r.Spec = tmp.(Requirement)
	}

	return err
}

// RequirementOperandDeclarator method. Declares type for the "spec" field
func RequirementOperandDeclarator (r RequirementOperand) { //Error checking not done;Case raw_value not done
	switch r.Type {
	case "raw_value":


	case "request":
		tmp := r.Spec
		r.Spec = tmp.(Hook)
	}
}


func main() {
	
}
