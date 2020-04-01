package common

import "github.com/pquerna/ffjson/ffjson"

//EnJson 序列化
//@param interface{} v
//@return []byte, error
func EnJson(v interface{}) ([]byte, error) {
	return ffjson.Marshal(&v)
}

//DeJson 反序列化
//@param []byte data
//@param interface{} v
//@return error
func DeJson(data []byte, v interface{}) error {
	if len(data) < 1 {
		return nil
	}
	return ffjson.Unmarshal(data, v)
}
