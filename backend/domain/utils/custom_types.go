package utils

import (
	"encoding/json"
	"net/url"
)

// Represent Dictonary type
type Dict map[string]interface{}

func (d Dict) Get(keys ...string) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if len(keys) == 1 {
		return d[keys[0]]
	}
	t, ok := d[keys[0]].(Dict)
	if !ok {
		return nil
	}
	return t.Get(keys[1:]...)
}

func (d *Dict) ToString() string {
	str, _ := json.Marshal(d)
	return string(str)
}

func (d *Dict) ToUrlString() string {
	data := url.Values{}
	for k, v := range *d {
		if v != "" {
			data.Set(k, v.(string))
		}
	}
	return data.Encode()
}

// Represent List type
type List []interface{}

type Any interface{}
