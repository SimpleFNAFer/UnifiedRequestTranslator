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
	`{
  "source": "[connection string will be here]",
  "modifiers": "[modifiers will be here]",
  "requirements":
  {
    "type":"and",
    "or_and":
    [
      {
        "type":"or",
        "or_and":
        [
          {
            "type":"requirement",
            "requirement":
            {
              "operator":"in",
              "field": {"child":null, "name":"Id"},
              "value": {"type": "raw_value", "spec": [10,11]}
            }
          },
          {
            "type":"requirement",
            "requirement":
            {
              "operator":"gt",
              "field": {"child":null, "name":"Price"},
              "value": {"type": "raw_value", "spec": 150}
            }
          }
        ]
      },
      {
        "type":"requirement",
        "requirement":
        {
          "operator":"eq",
          "field": {"child": {"child": null, "name": "Name"}, "name":"Fio"},
          "value": {"type": "raw_value", "spec": "Joakim"}
        }
      },
      {
        "type":"requirement",
        "requirement":
        {
          "operator":"eq",
          "field": {"child": {"child": null, "name": "Surname"}, "name":"Fio"},
          "value": {"type": "raw_value", "spec": "Broden"}
        }
      }
    ]
  },
  "fields":
  [
    {"child":null, "name":"id"},
    {"child":
    {"child":null, "name":"Name"},
      "name":"Fio"
    },
    {"child":
    {"child": null, "name": "Surname"},
      "name": "Fio"
    }
  ]
}`,
}

var expectedSQL = []string{
	"SELECT id, Fio.Name FROM [connection string will be here] WHERE (id = '10')",
	"SELECT id, Fio.Name, Fio.Surname FROM [connection string will be here] WHERE " +
		"(((Id IN ('10', '11')) OR (Price > 150)) AND (Fio.Name = 'Joakim') AND (Fio.Surname = 'Broden'))",
}

var expectedES = []string{
	`{"_source":["id","Fio.Name"],"query":{"term":{"id":10}}}`,
	`How to make an 'IN' equivalent in ES?`,
}

func ReceiveFromStr(s string) (tr.UnifiedRequest, error) {
	var ur tr.UnifiedRequest
	var e = json.Unmarshal([]byte(s), &ur)

	return ur, e
}

func TestToSql(t *testing.T) {
	var (
		received string
		ur       tr.UnifiedRequest
	)

	for i, s := range src {
		ur, _ = ReceiveFromStr(s)
		received = ur.UnifiedRequestToSql()

		if received != expectedSQL[i] {
			t.Error("for: ", src[i], " expected: ", expectedSQL[i], ", but received: ", received)
		}
	}
}

func TestToES(t *testing.T) {
	var (
		received []byte
		ur       tr.UnifiedRequest
	)

	for i, s := range src {
		ur, _ = ReceiveFromStr(s)
		received, _ = json.Marshal(ur.UnifiedRequestToES())

		if string(received) != expectedES[i] {
			t.Error("for: ", src[i], " expected: ", expectedES[i], ", but received: ", string(received))
		}
	}
}
