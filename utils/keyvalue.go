package utils

import (
	"fmt"
	// "strings"
	"reflect"
	"time"
	"bytes"
	"encoding/json"
	// "github.com/user/ibsi"
)

// func ifs(condition bool, a string, b string) string {
	// if condition {
		// return a
	// }
	
	// return b
// }

// https://code.tutsplus.com/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540
// by convention, Go interface names end with "er"
type IKeyValue interface {
	Append(src IKeyValue)
	AsDateTime(key string) time.Time
	AsInt(key string) int
	AsString(key string) string
	ContainsKey(name string) bool
	Count() int
	Get(key string) interface{}
	GetCollection() map[string]interface{}
	Json() string
	JsonPair(key string) string
	// JsonSort(sort string) string
	Set(key string, value interface{})
	// Sort(sort string)
	// SetEx(key string, value interface{}, overwrite bool)
	Type(key string) string
}

type KeyValue struct {
	types map[string]string
	values map[string]interface{}
	sort []string
}

func NewKeyValue() IKeyValue {
    return &KeyValue{
		types: make(map[string]string), 
		values: make(map[string]interface{}),
		sort: make([]string, 0),
	}
}

func NewKeyValueEx(init string) IKeyValue {
    return &KeyValue{
		types: make(map[string]string), 
		values: make(map[string]interface{}),
		sort: make([]string, 0),
	}
}

// func NewKeyValueEx(init string) IKeyValue {
    // return &KeyValue{make(map[string]string), make(map[string]interface{})}
// }

func (kv KeyValue) Append(src IKeyValue) {
	for n := range src.GetCollection() {
		kv.Set(n, src.Get(n))
	}
}

func (kv *KeyValue) Set(key string, value interface{}) {
	if !kv.ContainsKey(key) {
		kv.sort = append(kv.sort, key)
	}
	
	kv.types[key] = reflect.TypeOf(value).Name()
	kv.values[key] = value
}

// func (kv KeyValue) SetEx(key string, value interface{}, overwrite bool) {
	// kv.types[key] = reflect.TypeOf(value).Name()
	// kv.types[key] = reflect.TypeOf(value).String()
	// kv.values[key] = value
// }

func (kv KeyValue) AsBoolean(key string) bool {
  return kv.values[key].(bool)
}

func (kv KeyValue) AsDateTime(key string) time.Time {
  return kv.values[key].(time.Time)
}

func (kv KeyValue) AsInt(key string) int {
  return kv.values[key].(int)
}

func (kv KeyValue) AsString(key string) string {
  return kv.values[key].(string)
}

func (kv KeyValue) Type(key string) string {
  return kv.types[key]
}

func (kv KeyValue) Get(key string) interface{} {
  return kv.values[key]
}

func (kv KeyValue) JsonPair(key string) string {
	var s string
	v := kv.Get(key)	
	// t := kv.Type(key)
	
	// if t == "string" {
		// s = fmt.Sprintf("\"%v\"", v)
	// } else if t == "Time" {
		// s = fmt.Sprintf("\"%s\"", v.(time.Time).Format("2006-01-02T15-04-05.000"))
	// } else {
		var vs string
		json, err := json.MarshalIndent(v, "\t", "\t") // \t is formatted, remove it to minify
		if err != nil {
			vs = "{}"
		} else {
			vs = string(json[:])
		}
	
		s = fmt.Sprintf("%s", vs)
		// s = fmt.Sprintf("%v", v)
	// }
	
	// return fmt.Sprintf(string('"%s":%v'), key, v)
	return fmt.Sprintf("\"%s\": %s", key, s)
}

// func (kv KeyValue) Json() string {
	// total := kv.Count()
	// if (total == 0) {
		// return "{}"
	// }
	
	// var k = 1;
	// var json bytes.Buffer;
	
	// json.WriteString("{\n")
	// for i := range kv.GetCollection() {		
		// json.WriteString(fmt.Sprintf("\t%s%s\n", kv.JsonPair(i), Ifs(k < total, ",", "")))
		// k++
	// }
	// json.WriteString("}\n")
	
	// return json.String()
// }

func (kv KeyValue) Json() string {
	// if (sort == "") {
		// return "{}"
	// }
	
	var k = 1;
	var json bytes.Buffer;
	// names := strings.Split(kv.sort, ",")
	total := len(kv.sort)
	
	json.WriteString("{\n")
	// for i := range kv.GetCollection() {		
	for i := 0; i < len(kv.sort); i++ {		
		json.WriteString(fmt.Sprintf("\t%s%s\n", kv.JsonPair(kv.sort[i]), Ifs(k < total, ",", "")))
		k++
	}
	json.WriteString("}\n")
	
	return json.String()
}

func (kv KeyValue) ContainsKey(name string) bool {
	// if _, ok := dict["foo"]; ok {
		// return true
	// }	
	_, yes := kv.values[name]
	return yes
}

func (kv KeyValue) Count() int {
  return len(kv.GetCollection())
}

func (kv KeyValue) GetCollection() map[string]interface{} {
  return kv.values
}

// func (kv *KeyValue) Sort(sort string) {
  // kv.sort = sort
// }
