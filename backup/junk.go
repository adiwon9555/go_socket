package backup

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Name string
	Data interface{}
}

func main1() {
	rawMsg := []byte(`{"name":"channel add",` +
		`"data":{"name":"Hardware Support"}}`)
	var stdata Message
	json.Unmarshal(rawMsg, &stdata)
	fmt.Println(stdata)
}
