package common

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func StringHasPrefix(input, prefix string) bool {
	return strings.HasPrefix(input, prefix)
}
func StringTrimPrefix(input, prefix string) string {
	return strings.TrimPrefix(input, prefix)
}
func StringContains(input, substring string) bool {
	return strings.Contains(input, substring)
}
func StringToUpper(input string) string {
	return strings.ToUpper(input)
}
func StringToLower(input string) string {
	return strings.ToLower(input)
}

func IntToString(input int) string {
	return strconv.Itoa(input)
}

func StringToFloat64(input string) float64 {
	float, _ := strconv.ParseFloat(input, 64)

	return float
}

func Int64ToString(input int64) string {
	return strconv.FormatInt(input, 10)
}

func StringToInt64(input string) int64 {
	int, _ := strconv.ParseInt(input, 10, 64)

	return int
}

func StringToInt(input string) int {
	int, _ := strconv.Atoi(input)

	return int
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func StringReplace(s, old, new string) string {
	return strings.Replace(s, old, new, 0)
}
func StringDataCompareUpper(haystack []string, needle string) bool {
	for x := range haystack {
		if StringToUpper(haystack[x]) == StringToUpper(needle) {
			return true
		}
	}
	return false
}
func JSONDecode(data []byte, to interface{}) error {
	if !StringContains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}
func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}


