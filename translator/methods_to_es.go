package translator

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
	Should  []Query `json:"should,omitempty"`
	Must    []Query `json:"must,omitempty"`
	MustNot []Query `json:"must_not,omitempty"`
}

func (u UnifiedRequest) UnifiedRequestToES() ESRequest {
	var tmpRequest ESRequest

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
		temp.Bool = &Bool{}
		for _, i := range r.OrAnd {
			temp.Bool.Should = append(temp.Bool.Should, i.ToESRequirementExpression())
		}
		return temp

	case "and":
		temp.Bool = &Bool{}
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
		temp.Terms[r.Field.TransferToString()] = r.Value.Spec.([]interface{})

	case "not in":
		temp.Bool = &Bool{
			MustNot: []Query{
				{
					Terms: map[string][]interface{}{r.Field.TransferToString(): r.Value.Spec.([]interface{})},
				},
			},
		}

	case "ne":
		temp.Bool = &Bool{
			MustNot: []Query{
				{
					Term: map[string]interface{}{r.Field.TransferToString(): r.Value.Spec},
				},
			},
		}

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
