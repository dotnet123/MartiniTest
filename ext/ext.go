package ext

import (
	"encoding/json"
	"strconv"
)


type NullableString struct {
	Value   string
	HasValue bool
}

func (n *NullableString) Set(value string) error {
	n.HasValue = true
	n.Value = value
	return nil
}
func (n *NullableString) MarshalJSON() ([]byte, error) {
	if n.HasValue {
		//sn := string(n.Data)
		return json.Marshal(n.Value)
	}

	var i interface{}
	return json.Marshal(i)

}

func (ns *NullableString) UnmarshalJSON(b []byte) error {

	str := string(b)
	unquoted,err:= strconv.Unquote(str)
	if err != nil{
		return err
	}
	// Special case when we encounter `null`, modify it to the empty string
	if str == "null" || str == "" {

	} else {

		*ns = NullableString{ unquoted,true}
	}

	return nil


}
