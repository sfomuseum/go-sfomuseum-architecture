package properties

import (
	"github.com/tidwall/gjson"
)

func Concordances(body []byte) map[string]string {

	rsp := gjson.GetBytes(body, "properties.wof:concordances")

	if !rsp.Exists() {
		return nil
	}

	concordances := make(map[string]string)

	for k, v := range rsp.Map() {
		concordances[k] = v.String()
	}

	return concordances
}
