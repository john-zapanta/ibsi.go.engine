package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func Contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}

func StrToInt(expr string) int {
	v, _ := strconv.Atoi(expr)
	return v
}

func StrToInt64(expr string) int64 {
	v, _ := strconv.Atoi(expr)
	return int64(v)
}

func Ifs(condition bool, a string, b string) string {
	if condition {
		return a
	}
	
	return b
}

func Ifx(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	
	return b
}

func Ifempty(input string, def string) string {
	if input == "" {
		return def
	}
	
	return input
}

func CreateMapFromSlices(keys []string, values []interface{}) map[string]interface{} {
	kv := make(map[string]interface{})
	for i := 0; i < len(keys); i++ {
		kv[keys[i]] = values[i]
	}
	
	return kv
}

func SplitName(name string, defaultConnection string) (string, string) {
	parts := strings.Split(name, ".")
	
	var cname, command string = defaultConnection, ""
	if len(parts) == 2 {
		cname = parts[0]
		command = parts[1]
	} else {
		command = parts[0]
	}
	
	return cname, command
}

// thanks to: https://stackoverflow.com/questions/22491876/convert-byte-array-uint8-to-float64-in-golang
// func Float64FromBytes(bytes []byte) float64 {
    // bits := binary.LittleEndian.Uint64(bytes)
    // float := math.Float64frombits(bits)
    // return float
// }

// func Float64FromBytes2(bytes []uint8) float64 {
    // bits := binary.LittleEndian.Uint64(bytes)
    // float := math.Float64frombits(bits)
    // return float
// }

func Println(args ...interface{}) {
	fmt.Println(args...)
}