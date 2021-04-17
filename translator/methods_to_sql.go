package translator

import "fmt"

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
