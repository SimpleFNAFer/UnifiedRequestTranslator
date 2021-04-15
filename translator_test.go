package test

import (
	tr "UnifiedRequestTranslator/translator"
	"encoding/json"
	"testing"
)

var src = []string{
	`{
	"source": "[connection string will be here]",
	"modifiers": "[modifiers will be here]",
	"requirements":
	{
		"type":"requirement",
		"requirement":
		{
			"operator":"eq",
			"field":{"child":null, "name":"id"},
			"value":{"type":"raw_value", "spec":10}
		}
	},
	"fields":
	[
		{"child":null, "name":"id"},
		{"child":
			{"child":null, "name":"Name"},
		"name":"Fio"
		}
	]
}`,
}

var expectedSQL = []string{
	"SELECT id, Fio.Name FROM [connection string will be here] WHERE (id = '10')",
}

func RecieveFromStr(s string) (tr.UnifiedRequest, error) {
	var ur tr.UnifiedRequest
	var e = json.Unmarshal([]byte(s), &ur)

	return ur, e
}

func TestToSql(t *testing.T) {
	var (
		recieved string
		ur       tr.UnifiedRequest
	)

	for i, s := range src {
		ur, _ = RecieveFromStr(s)
		recieved = ur.UnifiedRequestToSql()

		if recieved != expectedSQL[i] {
			t.Error("for: ", src, " expected: ", expectedSQL[i], ", but received: ", recieved)
		}
	}
}
