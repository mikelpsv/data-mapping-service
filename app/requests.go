package app

import (
	"net/http"
	"strconv"
)

func GetSimpleValue(r *http.Request, paramName string) (value string, exists bool) {
	val, exists := r.URL.Query()[paramName]
	if exists {
		return val[0], exists
	} else {
		return "", exists
	}

}

func GetSimpleValueAsInt(r *http.Request, paramName string) (value int64, exists bool, err error) {
	val_map, exists := r.URL.Query()[paramName]
	if exists {
		value, err = strconv.ParseInt(val_map[0], 10, 64)
	} else {
		value = 0
	}
	return value, exists, err
}
