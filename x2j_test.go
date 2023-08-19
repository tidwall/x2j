package x2j

import (
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	xmldata := `
		<myfancyxml>
			<taco id="facts-are-facts">
				<burrito>Yum!</burrito>
				<burrito>Yummer!</burrito>
			</taco>
		</myfancyxml>
	`
	jsondata, _ := Convert([]byte(xmldata))
	if !json.Valid(jsondata) {
		t.Fail()
	}
	expect := `{"name":"myfancyxml","children":[{"name":"taco","attrs":` +
		`{"id":"facts-are-facts"},"children":[{"name":"burrito","children":` +
		`["Yum!"]},{"name":"burrito","children":["Yummer!"]}]}]}`
	if string(jsondata) != expect {
		t.Fail()
	}
}
